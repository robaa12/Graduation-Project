import { PlansService } from './../plans/plans.service';
import { CategoryService } from './../category/category.service';
import { CreateStoreThemeDto } from './dto/create-store-theme.dto';
import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { CreateStoreDto } from './dto/create-store.dto';
import { UpdateStoreDto } from './dto/update-store.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { EmailService } from 'src/shared/services/email/email.service';
import { Repository } from 'typeorm';
import { Store } from './entities/store.entity';
import { UserService } from 'src/user/user.service';
import { InjectModel } from '@nestjs/mongoose';
import { StoreThemeSchema } from './entities/store-theme.entity';
import { Model } from 'mongoose';
import { UpdateStoreThemeDto } from './dto/update-store-theme.dto';
import { StoreGallery } from './entities/user-gallery.entity';

@Injectable()
export class StoreService {
  constructor(
    @InjectRepository(Store) private storeRepository: Repository<Store>,
    @InjectRepository(StoreGallery)
    private storeGalleryRepository: Repository<StoreGallery>,
    private CategoryService: CategoryService,
    private readonly UserService: UserService,
    @InjectModel('StoreTheme') private storeThemeModel: Model<StoreThemeSchema>,
  ) {}

  async createStore(createStoreDto: CreateStoreDto): Promise<Store> {
    const user = await this.UserService.findOne(createStoreDto.user_id);
    if (!user) {
      throw new NotFoundException('User not found');
    }
    const category = await this.CategoryService.findOne(
      createStoreDto.category_id,
    );
    console.log(user);

    if (user.stores.length + 1 > user.plan.num_of_stores) {
      throw new BadRequestException(
        'You have reached the maximum number of stores allowed for your plan',
      );
    }
    let slug = await this.generateStoreSlug(createStoreDto.store_name);

    const store = this.storeRepository.create({
      ...createStoreDto,
      user,
      category,
      slug,
    });
    return await this.storeRepository.save(store);
  }

  async findAll(): Promise<Store[]> {
    return await this.storeRepository.find({ relations: ['category', 'user'] });
  }

  async deleteStore(id: number): Promise<void> {
    const store = await this.storeRepository.findOne({ where: { id } });
    if (!store) {
      throw new NotFoundException('Store not found');
    }
    await this.storeRepository.delete(id);
  }

  async findStoreBySlug(slug: string): Promise<Store> {
    const store = await this.storeRepository.findOne({
      where: { slug },
      relations: ['category', 'user'],
    });
    if (!store) {
      throw new NotFoundException('Store not found');
    }
    return store;
  }

  async findOne(id: number): Promise<Store> {
    const store = await this.storeRepository.findOne({
      where: { id },
      relations: ['category', 'user'],
    });
    if (!store) {
      throw new NotFoundException('Store not found');
    }
    return store;
  }

  async findAllStoresByUserId(userId: number): Promise<Store[]> {
    const user = await this.UserService.findOne(userId);
    if (!user) {
      throw new NotFoundException('User not found');
    }
    return await this.storeRepository.find({
      where: {
        user: {
          id: userId,
        },
      },
      relations: ['category'],
    });
  }

  async createStoreTheme(CreateStoreThemeDto: CreateStoreThemeDto) {
    const store = await this.storeRepository.findOne({
      where: { id: CreateStoreThemeDto.storeId },
    });
    if (!store) {
      throw new NotFoundException('Store not found');
    }
    console.log(CreateStoreThemeDto.theme.selectedTheme.id);

    let existingTheme = await this.storeThemeModel.findOne({
      storeId: CreateStoreThemeDto.storeId,
      'theme.selectedTheme.id': CreateStoreThemeDto.theme.selectedTheme.id,
    });
    if (existingTheme) {
      existingTheme = await this.storeThemeModel.findOneAndUpdate(
        { _id: existingTheme._id },
        {
          theme: CreateStoreThemeDto.theme,
          isActive: CreateStoreThemeDto.isActive,
        },
        { new: true, runValidators: true },
      );

      if (existingTheme.isActive) {
        await this.storeThemeModel.updateMany(
          {
            storeId: CreateStoreThemeDto.storeId,
            _id: { $ne: existingTheme._id },
          },
          { isActive: false },
        );
      }
      return existingTheme;
    }
    const storeTheme = await this.storeThemeModel.create(CreateStoreThemeDto);
    if (CreateStoreThemeDto.isActive) {
      await this.storeThemeModel.updateMany(
        { storeId: CreateStoreThemeDto.storeId, _id: { $ne: storeTheme._id } },
        { isActive: false },
      );
    }
    return storeTheme;
  }

  async findStoreThemes(storeId: number) {
    return await this.storeThemeModel.find({ storeId });
  }
  async fincStoreActiveTheme(storeId: number) {
    const storeTheme = await this.storeThemeModel.findOne({
      storeId,
      isActive: true,
    });
    return storeTheme;
  }

  async findStoreActiveThemeByStoreSlug(storeSlug: string) {
    const store = await this.storeRepository.findOne({
      where: { slug: storeSlug },
    });
    if (!store) {
      throw new NotFoundException('Store not found');
    }
    const storeTheme = await this.storeThemeModel.findOne({
      storeId: store.id,
      isActive: true,
    });
    if (!storeTheme) {
      throw new NotFoundException('Store theme not found');
    }
    return storeTheme;
  }

  async findStoreThemesByStoreId(id: string) {
    return await this.storeThemeModel.findOne({ _id: id });
  }

  async removeStoreTheme(id: string) {
    return await this.storeThemeModel.findOneAndDelete({ _id: id });
  }

  async generateStoreSlug(storeName: string): Promise<string> {
    let slug = storeName.toLocaleLowerCase().replace(/ /g, '-');
    let existingStore = await this.storeRepository.findOne({ where: { slug } });

    let counter = 1;
    while (existingStore) {
      slug = `${storeName.toLocaleLowerCase().replace(/ /g, '-')}-${counter}`;
      existingStore = await this.storeRepository.findOne({ where: { slug } });
      counter++;
    }

    return slug;
  }

  async addPhotoToGallery(
    storeId: number,
    imageUrl: string,
  ): Promise<StoreGallery> {
    const store = await this.storeRepository.findOneBy({ id: storeId });
    if (!store) {
      throw new NotFoundException('User not found');
    }
    const galleryItem = this.storeGalleryRepository.create({
      store,
      imageUrl,
    });
    return await this.storeGalleryRepository.save(galleryItem);
  }

  async getStoreGallery(storeId: number): Promise<StoreGallery[]> {
    const gallery = await this.storeGalleryRepository.find({
      where: { store: { id: storeId } },
    });
    if (!gallery) {
      throw new NotFoundException('Gallery not found for this user');
    }
    return gallery;
  }

  async deletePhotoFromGallery(image_id: number): Promise<void> {
    const galleryItem = await this.storeGalleryRepository.findOneBy({
      id: image_id,
    });
    if (!galleryItem) {
      throw new NotFoundException('Gallery item not found');
    }
    await this.storeGalleryRepository.remove(galleryItem);
  }

  async addPhotosToGallery(storeId: number, imageUrls: string[]) {
    const store = await this.storeRepository.findOneBy({ id: storeId });
    if (!store) {
      throw new NotFoundException('Store not found');
    }

    const galleryItems = imageUrls.map((imageUrl) =>
      this.storeGalleryRepository.create({
        store,
        imageUrl,
      }),
    );

    const images = await this.storeGalleryRepository.save(galleryItems);
    return images.map((image) => {
      return {
        id: image.id,
        imageUrl: image.imageUrl,
      };
    });
  }
}

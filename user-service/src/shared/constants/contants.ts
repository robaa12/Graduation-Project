import { CreateCategoryDto } from "src/category/dto/create-category.dto";
import { CreatePlanDto } from "src/plans/dto/create-plan.dto";

export const defaultPlans: CreatePlanDto[] = [
    {
        name: 'Free Plan',
        price: 0,
        description: 'Free plan with basic features',
        num_of_stores: 1,
      },
      {
        name: 'Basic Plan',
        price: 500,
        description: 'Basic plan with essential features',
        num_of_stores: 2,
      },
      {
        name: 'Pro Plan',
        price: 1000,
        description: 'Pro plan with additional features',
        num_of_stores: 5,
      },
      {
        name: 'Premium Plan',
        price: 1500,
        description: 'Premium plan with all features included',
        num_of_stores: 10,
      },
];


export const defaultCategories:CreateCategoryDto[] = [
  {
    "name": "Electronics",
    "description": "Devices and gadgets such as phones, laptops, cameras, and accessories."
  },
  {
    "name": "Fashion",
    "description": "Clothing, footwear, and accessories for men, women, and children."
  },
  {
    "name": "Home & Furniture",
    "description": "Furniture, home décor, lighting, and kitchenware."
  },
  {
    "name": "Health & Beauty",
    "description": "Cosmetics, skincare, personal care, and wellness products."
  },
  {
    "name": "Groceries",
    "description": "Everyday food items, beverages, and household essentials."
  },
  {
    "name": "Sports & Outdoors",
    "description": "Sportswear, fitness gear, camping, and outdoor equipment."
  },
  {
    "name": "Toys & Games",
    "description": "Toys, board games, and educational items for children."
  },
  {
    "name": "Automotive",
    "description": "Car accessories, tools, and maintenance supplies."
  },
  {
    "name": "Books & Stationery",
    "description": "Books, notebooks, art supplies, and office materials."
  },
  {
    "name": "Pets",
    "description": "Pet food, accessories, and care products for animals."
  }
]

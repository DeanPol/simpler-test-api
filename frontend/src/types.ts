export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
}

export interface ProductCount {
  total: number;
  limit: number;
}

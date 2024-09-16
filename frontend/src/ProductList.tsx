import React, { useState, useEffect } from "react";
import axios from "axios";
import { Product, ProductCount } from "./types";

const ProductList: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(10);
  const [totalProducts, setTotalProducts] = useState<number>(0);

  useEffect(() => {
    const fetchProductCount = async () => {
      try {
        const response = await axios.get<ProductCount>(
          "http://localhost:8080/products/total"
        );
        setLimit(response.data.limit);
        setTotalProducts(response.data.total);
      } catch (err) {
        setError("Error fetching product count");
        console.error("Error fetching product count", err);
      }
    };

    const fetchProducts = async () => {
      try {
        const response = await axios.get<Product[]>(
          "http://localhost:8080/products",
          {
            params: {
              _page: page,
              _limit: limit,
            },
          }
        );
        setProducts(response.data);
      } catch (err) {
        setError("Error fetching products");
        console.error("Error fetching products", err);
      }
    };

    fetchProductCount();
    fetchProducts();
  }, [page, limit]);

  // Calculate total number of pages
  const totalPages = Math.ceil(totalProducts / limit);

  // Handlers for pagination buttons
  const handlePrevious = () => {
    if (page > 1) setPage(page - 1);
  };

  const handleNext = () => {
    if (page < totalPages) setPage(page + 1);
  };

  return (
    <div>
      <h1>Product List</h1>
      {error && <p>{error}</p>}
      <ul>
        {products.map((product) => (
          <li key={product.id}>
            {product.name} - ${product.price}
          </li>
        ))}
      </ul>
      <div>
        <button onClick={handlePrevious} disabled={page === 1}>
          Previous
        </button>
        <button onClick={handleNext} disabled={page === totalPages}>
          Next
        </button>
      </div>
    </div>
  );
};

export default ProductList;

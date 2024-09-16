import React, { useState, useEffect } from "react";
import axios from "axios";
import { Product } from "./types";

const ProductList: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [limit, setLimit] = useState<number>(1); // Start with 0, updated from backend
  const [page, setPage] = useState<number>(0); // Offset value
  const [totalProducts, setTotalProducts] = useState<number>(0);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await axios.get<{
          total: number;
          limit: number;
          offset: number;
          products: Product[];
        }>("localhost:8080/products", {
          params: {
            limit: limit > 0 ? limit : undefined, // Only include limit if it's set
            offset: page,
          },
        });

        // Set data from API response
        setProducts(response.data.products);
        setTotalProducts(response.data.total);

        // Set limit and page based on backend response
        if (limit === 0) setLimit(response.data.limit); // Only set the limit if it's not yet set
        setPage(response.data.offset);
      } catch (err) {
        setError("Error fetching products");
        console.error("Error fetching products", err);
      }
    };

    fetchProducts();
  }, [page, limit]); // Limit and page are dependencies

  // Calculate total number of pages
  const totalPages = Math.ceil(totalProducts / limit);

  // Handlers for pagination buttons
  const handlePrevious = () => {
    if (page > 0) setPage(page - 1);
  };

  const handleNext = () => {
    if (page < totalPages - 1) setPage(page + 1);
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
        <button onClick={handlePrevious} disabled={page === 0}>
          Previous
        </button>
        <button onClick={handleNext} disabled={page >= totalPages - 1}>
          Next
        </button>
      </div>
    </div>
  );
};

export default ProductList;

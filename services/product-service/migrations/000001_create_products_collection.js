db.createCollection("products");

db.products.createIndex({ name: 1 }, { unique: true });
db.products.createIndex({ category: 1 });
db.products.createIndex({ created_at: -1 });

db.products.createIndex(
  { name: "text", description: "text" },
  {
    weights: {
      name: 10,
      description: 5,
    },
    name: "products_text_search",
  }
);

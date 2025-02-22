#!/bin/bash

# Project Name
PROJECT_NAME="vite-react-app"

# Create Vite Project
echo "🚀 Creating Vite project..."
npm create vite@latest $PROJECT_NAME -- --template react

# Change into project directory
cd $PROJECT_NAME

# Install dependencies
echo "📦 Installing dependencies..."
npm install
npm install -D tailwindcss postcss autoprefixer
npm install react-router-dom

# Initialize Tailwind CSS
echo "🎨 Setting up Tailwind CSS..."
npx tailwindcss init -p

# Configure Tailwind
echo "🔧 Configuring Tailwind..."
cat <<EOT > tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
};
EOT

# Add Tailwind to CSS
echo "📝 Updating CSS..."
cat <<EOT > src/index.css
@tailwind base;
@tailwind components;
@tailwind utilities;
EOT

# Create basic Router setup
echo "🛤 Setting up React Router..."
cat <<EOT > src/main.jsx
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")).render(
  <BrowserRouter>
    <App />
  </BrowserRouter>
);
EOT

cat <<EOT > src/App.jsx
import { Routes, Route, Link } from "react-router-dom";

export default function App() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <h1 className="text-3xl font-bold text-blue-600">Welcome to Vite + React + Tailwind + React Router</h1>
      <nav className="mt-4">
        <Link className="text-blue-500 hover:underline mx-2" to="/">Home</Link>
        <Link className="text-blue-500 hover:underline mx-2" to="/about">About</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </div>
  );
}

function Home() {
  return <h2 className="text-xl">🏡 Home Page</h2>;
}

function About() {
  return <h2 className="text-xl">ℹ️ About Page</h2>;
}
EOT

# Done!
echo "✅ Setup complete!"
echo "🔄 Run the project with:"
echo "cd $PROJECT_NAME && npm run dev"

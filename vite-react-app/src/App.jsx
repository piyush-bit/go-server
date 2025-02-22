import { Routes, Route, Link, useNavigate } from "react-router-dom";

export default function App() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <Routes>
        <Route path="/" element={<AuthPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </div>
  );
}

import React, { useState } from 'react';
import { Github, Mail, LogIn } from 'lucide-react';
import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

const AuthPage = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [companyData, setCompanyData] = useState(null);

  const useQuery = () => {
    return new URLSearchParams(useLocation().search);
  };

  const query = useQuery();
  const id = query.get('id');

  const navigate = useNavigate();

  const handleSignup = async () => {
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    formData.append('name', name);
    if(companyData != null) {
      formData.append('app_id', companyData.id);
    }

    try {
      const response = await fetch('http://localhost:8080/api/v1/signup', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      console.log('Signup successful:', data);

      if(companyData != null) {
        // Redirect to company login page
        window.location.href = companyData.callback_url+"?email="+email+"&name="+name+'&token_id='+data.data.token_id;
      }
      // store token in local storage 
      localStorage.setItem('token', data.data.token);
      navigate('/dashboard');
      // Handle successful signup
    } catch (error) {
      console.error('Error during signup:', error);
    }
  };

  const handleLogin = async () => {
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);
    if(companyData != null) {
      formData.append('app_id', companyData.id);
    }
    try {
      const response = await fetch('http://localhost:8080/api/v1/login', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      console.log('Login successful:', data);

      if(companyData != null) {
        // Redirect to company login page
        window.location.href = companyData.callback_url+"?email="+email+"&name="+name+'&token_id='+data.data.token_id;
      }
      // store token in local storage 
      localStorage.setItem('token', data.data.token);
      navigate('/dashboard');
      
    } catch (error) {
      console.error('Error during login:', error);
    }
  };

  useEffect(() => {
    if(id == null) {
      return;
    }
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/v1/app/get/'+id);
        const data = await response.json();
        console.log(data);
        setCompanyData(data.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);



  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle form submission
    console.log('Form submitted:', { email, password, name });
  };

  return (
    <div className="min-h-screen flex items-center justify-center ">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-lg">
        {/* Company Branding */}
        {<div className="text-center space-y-2">
          <div className="flex items-center justify-center space-x-3">
            {/* Placeholder for company logo */}
            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center flex-shrink-0">
              <span className="text-xl font-bold text-blue-600">Logo</span>
            </div>
            <h1 className="text-2xl font-bold text-gray-900">{companyData?.name ?? "Go SSO"}</h1>
          </div>
          <p className="text-sm text-gray-500">Single Sign-On Portal</p>
        </div>}

        <div className="border-t border-gray-200 pt-8">
          {/* Header */}
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900">
              {isLogin ? 'Welcome back' : 'Create account'}
            </h2>
            <p className="mt-2 text-sm text-gray-600">
              {isLogin ? "Don't have an account? " : "Already have an account? "}
              <button
                onClick={() => setIsLogin(!isLogin)}
                className="text-blue-600 hover:text-blue-500 font-medium transition-colors"
              >
                {isLogin ? 'Sign up' : 'Log in'}
              </button>
            </p>
          </div>
        </div>

        {/* SSO Buttons */}
        <div className="space-y-3 hidden">
          <button className="w-full flex items-center justify-center gap-3 px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors">
            <Mail className="h-5 w-5" />
            Continue with Google
          </button>
          <button className="w-full flex items-center justify-center gap-3 px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors">
            <Github className="h-5 w-5" />
            Continue with GitHub
          </button>
        </div>

        <div className="relative hidden">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-white text-gray-500">Or continue with</span>
          </div>
        </div>

        {/* Email Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
          {!isLogin && (
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Name
              </label>
              <input
                id="name"
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          )}
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">
              Email address
            </label>
            <input
              id="email"
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              id="password"
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          {isLogin && (
            <div className="flex items-center justify-end">
              <button className="text-sm text-blue-600 hover:text-blue-500">
                Forgot password?
              </button>
            </div>
          )}
          <button
            onClick={isLogin ? handleLogin : handleSignup}
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
          >
            <LogIn className="h-5 w-5 mr-2" />
            {isLogin ? 'Sign in' : 'Create account'}
          </button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center text-sm text-gray-500">
          By continuing, you agree to Company's{' '}
          <a href="#" className="text-blue-600 hover:text-blue-500">
            Terms of Service
          </a>{' '}
          and{' '}
          <a href="#" className="text-blue-600 hover:text-blue-500">
            Privacy Policy
          </a>
        </div>
      </div>
    </div>
  );
};

function Dashboard() {
  return <h2 className="text-xl">ℹ️ Dashboard </h2>;
}

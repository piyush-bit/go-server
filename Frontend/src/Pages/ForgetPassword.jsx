import React, { useState } from 'react';
import { Link } from 'react-router-dom';

const ForgotPasswordPage = () => {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const BACKEND_URI = import.meta.env.VITE_BACKEND_URI??"";

  const validateEmail = (email) => {
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setMessage('');
    
    if (!validateEmail(email)) {
      setError('Please enter a valid email address');
      return;
    }
    
    setLoading(true);
    
    try {
      // Send password reset request to the backend
      const formData = new FormData();
      formData.append('email', email);
      const response = await fetch(BACKEND_URI+'/api/v1/forget-password', {
        method: 'POST',
        body: formData,
      });
      
      if (response.ok) {
        setMessage('Password reset instructions have been sent to your email');
        setEmail('');
      } else {
        const data = await response.json();
        setError(data.message || 'Failed to send reset email. Please try again.');
      }
    } catch (err) {
      setError('An error occurred. Please try again later.');
      console.error('Forgot password error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
        <h2 className="mb-6 text-2xl font-bold text-center text-gray-800">Forgot Your Password?</h2>
        <p className="mb-6 text-center text-gray-600">
          Enter your email address below and we'll send you instructions to reset your password.
        </p>
        
        {message && (
          <div className="p-3 mb-4 text-sm text-green-700 bg-green-100 rounded">
            {message}
          </div>
        )}
        
        {error && (
          <div className="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block mb-2 text-sm font-medium text-gray-700" htmlFor="email">
              Email Address
            </label>
            <input
              type="email"
              id="email"
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your-email@example.com"
              required
            />
          </div>
          
          <button
            type="submit"
            className="w-full px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
            disabled={loading}
          >
            {loading ? 'Sending...' : 'Send Reset Link'}
          </button>
        </form>
        
        <div className="mt-6 text-center">
          <div className="mb-2">
            <Link to={"/"}>
                <p className="text-sm text-blue-600 hover:underline">
                Back to Login
                </p>
            </Link>
          </div>
          <div>
          <Link to={"/"}>
                <p className="text-sm text-blue-600 hover:underline">
                Create Account
                </p>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ForgotPasswordPage;
import React, { useState, useEffect } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { Eye, EyeOff } from 'lucide-react';

const PasswordResetForm = () => {
  const [searchParams] = useSearchParams();
  const [email, setEmail] = useState('');
  const [token, setToken] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [confirmError, setConfirmError] = useState('');

  useEffect(() => {
    // Get email and token from URL parameters
    const emailParam = searchParams.get('email');
    const tokenParam = searchParams.get('token');
    
    if (emailParam) setEmail(emailParam);
    if (tokenParam) setToken(tokenParam);
  }, [searchParams]);


  const BACKEND_URI = import.meta.env.VITE_BACKEND_URI??"";
  // Check if confirm password matches in real-time
  useEffect(() => {
    if (confirmPassword && password !== confirmPassword) {
      setConfirmError('Passwords do not match');
    } else {
      setConfirmError('');
    }
  }, [password, confirmPassword]);

  const validatePassword = () => {
    if (password.length < 8) {
      setError('Password must be at least 8 characters long');
      return false;
    }
    
    if (password !== confirmPassword) {
      setConfirmError('Passwords do not match');
      return false;
    }
    
    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setMessage('');
    
    if (!validatePassword()) return;
    
    setLoading(true);
    
    try {
      const formData = new FormData();
      formData.append('password', password);
      formData.append('email', email);
      formData.append('token', token);
      
      const response = await fetch(BACKEND_URI+`/api/v1/reset-password`, {
        method: 'POST',
        body: formData,
      });
      
      const data = await response.json();
      
      if (response.ok) {
        setMessage('Password reset successful. You can now login with your new password.');
        setPassword('');
        setConfirmPassword('');
      } else {
        setError(data.message || 'Failed to reset password. Please try again.');
      }
    } catch (err) {
      setError('An error occurred. Please try again later.');
      console.error('Reset password error:', err);
    } finally {
      setLoading(false);
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const toggleConfirmPasswordVisibility = () => {
    setShowConfirmPassword(!showConfirmPassword);
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
        <h2 className="mb-6 text-2xl font-bold text-center text-gray-800">Reset Your Password</h2>
        
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
              Email
            </label>
            <input
              type="email"
              id="email"
              className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-100"
              value={email}
              disabled
            />
          </div>
          
          <div className="mb-4">
            <label className="block mb-2 text-sm font-medium text-gray-700" htmlFor="password">
              New Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                id="password"
                className="w-full px-3 py-2 border border-gray-300 rounded-md"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
              />
              <button 
                type="button"
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-600"
                onClick={togglePasswordVisibility}
              >
                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </div>
            <p className="mt-1 text-xs text-gray-500">Password must be at least 8 characters long</p>
          </div>
          
          <div className="mb-6">
            <label className="block mb-2 text-sm font-medium text-gray-700" htmlFor="confirmPassword">
              Confirm New Password
            </label>
            <div className="relative">
              <input
                type={showConfirmPassword ? "text" : "password"}
                id="confirmPassword"
                className={`w-full px-3 py-2 border rounded-md ${
                  confirmError ? "border-red-500" : "border-gray-300"
                }`}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
              <button 
                type="button"
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-600"
                onClick={toggleConfirmPasswordVisibility}
              >
                {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </div>
            {confirmError && (
              <p className="mt-1 text-xs text-red-500">{confirmError}</p>
            )}
          </div>
          
          <button
            type="submit"
            className="w-full px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
            disabled={loading || confirmError}
          >
            {loading ? 'Processing...' : 'Reset Password'}
          </button>
        </form>
        
        <div className="mt-4 text-center">
        <Link to={"/"}>
                <p className="text-sm text-blue-600 hover:underline">
                Back to Login
                </p>
            </Link>
        </div>
      </div>
    </div>
  );
};

export default PasswordResetForm;
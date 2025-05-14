import { useGoogleLogin } from '@react-oauth/google';
import { useState } from 'react';

const GoogleAuthButton = ({ onGoogleLogin }) => {
  const [isLoading, setIsLoading] = useState(false);

  const login = useGoogleLogin({
    onSuccess: async (tokenResponse) => {
      setIsLoading(true);
      try {
        console.log('Login Successful:', tokenResponse);
        const accessToken = tokenResponse.access_token;
        console.log('accessToken:', accessToken);
        if (onGoogleLogin) {
          await onGoogleLogin(accessToken);
        }
      } catch (error) {
        console.error('Error processing Google login:', error);
      } finally {
        setIsLoading(false);
      }
    },
    onError: error => {
      console.error('Login Failed:', error);
      setIsLoading(false);
    },
  });

  return (
    <button 
      onClick={() => login()} 
      disabled={isLoading}
      className=""
    >
      {isLoading ? 'Logging in...' : 'Continue with Google'}
    </button>
  );
};

export default GoogleAuthButton;
import React, { useState } from 'react';
import { useRouter } from 'next/router';

const LoginRegisterPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async () => {
    setLoading(true);
    setError(null);
  
    try {
      const response = await fetch(`http://localhost:8080/api/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });
  
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('token', data.token); // Set JWT token in localStorage
        router.push('/main');
      } else {
        const errorMessage = await response.text();
        setError(errorMessage || 'Login failed. Please try again.');
      }
    } catch (error) {
      setError('An error occurred. Please try again.');
    }
  
    setLoading(false);
  };
  
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <div className="bg-gray-300 rounded-2xl p-6 mb-4">
        <div className="flex flex-col items-center">
          <h1 className="text-3xl mb-4">Login Page</h1>
          <div className="mb-2">
            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="border text-black rounded-md px-4 py-2 w-72"
            />
          </div>
          <div className="mb-2">
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="border text-black rounded-md px-4 py-2 w-72"
            />
          </div>
          {error && <p style={{ color: 'red' }}>{error}</p>}
          <div>
            <button onClick={handleSubmit} disabled={loading} className="bg-blue-500 text-white px-4 py-2 rounded-md">
              Login
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginRegisterPage;

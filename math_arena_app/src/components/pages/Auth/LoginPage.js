import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './styles/Login.css';

function Login() {
    const [formData, setFormData] = useState({
        username: '',
        password: '',
    });

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const [isLoggedIn, setIsLoggedIn] = useState(false); 
    const navigate = useNavigate();

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            const response = await axios.post('http://localhost:8080/api/signin', formData);
            console.log('Login successful:', response.data);

            localStorage.setItem('authToken', response.data.authorization);  
            localStorage.setItem('userId', response.data.user_id);  

            setIsLoggedIn(true);  
        } catch (err) {
            console.error('Error during login:', err.response?.data || err.message);
            setError(err.response?.data?.message || 'Invalid username or password');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (isLoggedIn) {
            navigate('/home');
        }
    }, [isLoggedIn, navigate]);

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (token) {
            navigate('/home');
        }
    }, [navigate]);

    return (
        <div className="login-container">
            <h2 className="login-title">Log In</h2>

            {error && <p className="error-message">{error}</p>}

            {loading && <p className="loading-message">Logging in...</p>}

            <form className="login-form" onSubmit={handleSubmit}>
                <div className="input-group">
                    <label htmlFor="username" className="input-label">Username:</label>
                    <input
                        type="text"
                        id="username"
                        name="username"
                        value={formData.username}
                        onChange={handleChange}
                        required
                        className="input-field"
                        placeholder="Enter your username"
                    />
                </div>
                <div className="input-group">
                    <label htmlFor="password" className="input-label">Password:</label>
                    <input
                        type="password"
                        id="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        required
                        className="input-field"
                        placeholder="Enter your password"
                    />
                </div>
                <button type="submit" className="submit-button" disabled={loading}>
                    {loading ? 'Logging in...' : 'Log In'}
                </button>
            </form>

            <div className="redirect-message">
                <p>Don't have an account? <a href="/signup" className="redirect-button">Sign Up</a></p>
            </div>
        </div>
    );
}

export default Login;

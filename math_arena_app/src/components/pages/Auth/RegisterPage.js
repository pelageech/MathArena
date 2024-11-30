import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom'; 
import './styles/Register.css';


function Register() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
    });

    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(false);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    // Handle input change
    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

// Handle form submission with email validation
    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setSuccess(false);
        setLoading(true);

        // Email validation regex pattern (basic validation)
        const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;

        // Check if email is valid
        if (!emailRegex.test(formData.email)) {
            setError('Please enter a valid email address');
            setLoading(false);
            return;
        }

        try {
            // Sending POST request to the API endpoint
            const response = await axios.post('http://localhost:8080/api/signup', formData);
            console.log('Registration successful:', response.data);
            setSuccess(true);  // Set success message if registration is successful

            // Redirect to login page after 2 seconds
            setTimeout(() => {
                navigate('/login'); // Redirect to login page
            }, 2000);
        } catch (err) {
            console.error('Error during registration:', err.response?.data || err.message);
            setError(err.response?.data?.message || 'Something went wrong');  // Set error message if API fails
        } finally {
            setLoading(false);  // Stop loading spinner after the request is done
        }
    };








    return (
        <div className="register-container">
            <h2 className="register-title">Register</h2>

            {/* Success message */}
            {success && (
                <p className="success-message">Account created successfully! You can now log in.</p>
            )}

            {/* Error message */}
            {error && <p className="error-message">{error}</p>}

            {/* Loading state */}
            {loading && <p className="loading-message">Creating account...</p>}

            <form className="register-form" onSubmit={handleSubmit}>
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
                    <label htmlFor="email" className="input-label">Email:</label>
                    <input
                        type="email"
                        id="email"
                        name="email"
                        value={formData.email}
                        onChange={handleChange}
                        required
                        className="input-field"
                        placeholder="Enter your email"
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
                    {loading ? 'Creating...' : 'Register'}
                </button>
            </form>
        </div>
    );
}

export default Register;

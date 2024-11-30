import axios from 'axios';

// Create Game Session API
const createGameSession = async (userId) => {
    try {
        const response = await axios.post('http://localhost:8080/api/session/create', {
            user_id: userId,
        });
        return response.data;  // Return session details like session_id, time_left, and expression
    } catch (error) {
        console.error('Error creating game session:', error.response?.data || error.message);
        throw new Error('Failed to create game session');
    }
};

// Submit Answer API
const submitAnswer = async (sessionId, userId, answer) => {
    try {
        const response = await axios.post('http://localhost:8080/api/session/answer', {
            session_id: sessionId,
            user_id: userId,
            answer: answer,
        });
        return response.data;  // Return updated session details (expression, score, time_left)
    } catch (error) {
        console.error('Error submitting answer:', error.response?.data || error.message);
        throw new Error('Failed to submit answer');
    }
};

// Finish Game API
const finishGame = async (sessionId, userId) => {
    try {
        const response = await axios.post('http://localhost:8080/api/session/finish', {
            session_id: sessionId,
            user_id: userId,
        });
        return response.data;  // Return the game session status after finishing (if needed)
    } catch (error) {
        console.error('Error finishing game:', error.response?.data || error.message);
        throw new Error('Failed to finish game');
    }
};

export { createGameSession, submitAnswer, finishGame };

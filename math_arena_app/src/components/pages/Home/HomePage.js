import React, { useState, useEffect } from 'react';
import { createGameSession, submitAnswer, finishGame } from './Game/gameCardAPI';  // Import the combined API
import { useNavigate } from 'react-router-dom';
import './style/HomePage.css';

function HomePage() {
    const [gameData, setGameData] = useState(null);  // Store the created game session data
    const [answer, setAnswer] = useState(null);  // Store the user's answer
    const [gameSessionLoading, setGameSessionLoading] = useState(false);  // Loading state for session creation
    const [gameError, setGameError] = useState(null);  // Error state for game session creation
    const [gameOver, setGameOver] = useState(false);  // Flag to check if the game is over
    const [gameScore, setGameScore] = useState(0);  // Store the player's score
    const [timeLeft, setTimeLeft] = useState(0);  // Store the remaining time in seconds
    const [intervalId, setIntervalId] = useState(null);  // Store interval ID to clear later
    const [formattedTime, setFormattedTime] = useState(""); // Store formatted time
    const navigate = useNavigate();

    const token = localStorage.getItem('authToken');
    if (!token) {
        navigate('/login');
        return null;
    }

    // Convert nanoseconds to Minutes:Seconds
    const formatTime = (timeInNano) => {
        const timeInSeconds = timeInNano / 1000000000;  // Convert nanoseconds to seconds
        const minutes = Math.floor(timeInSeconds / 60);  // Get minutes
        const seconds = Math.floor(timeInSeconds % 60);  // Get remaining seconds
        return `${minutes}:${seconds < 10 ? '0' + seconds : seconds}`;  // Format as Minutes:Seconds
    };

    // Handle Game Session Creation
    const handleCreateGameSession = async () => {
        setGameSessionLoading(true);
        setGameError(null);
        try {
            const data = await createGameSession(4);  // Assuming user_id = 4
            setGameData(data);
            setFormattedTime(formatTime(data.time_left));  // Format the time from nanoseconds
            startTimer(data.time_left);  // Start the timer when the game session is created
        } catch (err) {
            setGameError('Failed to create game session');
        } finally {
            setGameSessionLoading(false);
        }
    };

    // Start the countdown timer for the game session
    const startTimer = (timeInNano) => {
        const timeInSeconds = timeInNano / 1000000000;  // Convert nanoseconds to seconds
        setTimeLeft(timeInSeconds);  // Set initial time left in seconds

        const id = setInterval(() => {
            setTimeLeft((prevTime) => {
                if (prevTime <= 0) {
                    clearInterval(id);
                    handleFinishGame();  // Finish the game if time runs out
                    return 0;
                }
                return prevTime - 1;  // Decrease the time by 1 second
            });
        }, 1000);
        setIntervalId(id);  // Store interval ID to clear it later
    };

    // Handle User's Answer Submission
    const handleSubmitAnswer = async () => {
        if (!answer) return;
        try {
            const data = await submitAnswer(gameData.session_id, 4, answer);
            setGameData(data);
            setGameScore(data.score);  // Update score if the answer is correct
            setAnswer(null);  // Clear the answer after submission
        } catch (err) {
            setGameError('Failed to submit answer');
        }
    };

    // Handle Finish Game
    const handleFinishGame = async () => {
        try {
            const result = await finishGame(gameData.session_id, 4);
            setGameOver(true);
            clearInterval(intervalId);  // Clear the timer interval
            alert(`Game Finished! Final Score: ${gameScore}`);
            navigate(0);  // Refresh the homepage
        } catch (err) {
            setGameError('Failed to finish game');
            // Delay the page refresh after showing the error message
            setTimeout(() => {
                navigate(0);  // Refresh the homepage after 2 seconds
            }, 2000);
        }
    };

    // Update the formatted time every second
    useEffect(() => {
        if (timeLeft > 0) {
            setFormattedTime(formatTime(timeLeft * 1000000000));  // Convert back to nanoseconds for formatting
        }
    }, [timeLeft]);

    return (
        <div className="home-container">
            <h1 className="home-title">Welcome to the MathArena</h1>

            {/* Button to create a new game session */}
            <button
                className="create-game-button"
                onClick={handleCreateGameSession}
                disabled={gameSessionLoading || gameData}
            >
                {gameSessionLoading ? 'Creating Game...' : 'Create Your Math Game'}
            </button>

            {/* Displaying the game session popup if the game session is created */}
            {gameData && !gameOver && (
                <div className="game-popup">
                    <h2>Math Expression: {gameData.expression}</h2>
                    <p>Time Left: {formattedTime}</p> {/* Show the formatted time */}
                    <p>Score: {gameScore}</p>

                    {/* Options to choose the answer */}
                    <input
                        type="number"
                        value={answer || ''}
                        onChange={(e) => setAnswer(Number(e.target.value))}
                        placeholder="Enter your answer"
                    />
                    <button onClick={handleSubmitAnswer}>Submit Answer</button>
                    
                    {/* Finish Game Button */}
                    <button onClick={handleFinishGame}>Finish Game</button>
                </div>
            )}

            {/* Game error message */}
            {gameError && <p className="error-message">{gameError}</p>}

            {/* Logout button */}
            <button
                className="logout-button"
                onClick={() => {
                    localStorage.removeItem('authToken');
                    navigate('/login');
                }}
            >
                Log out
            </button>
        </div>
    );
}

export default HomePage;

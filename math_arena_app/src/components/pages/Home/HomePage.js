import React, { useState, useEffect } from 'react';
import { createGameSession, submitAnswer, finishGame } from './Game/gameCardAPI';  
import { useNavigate } from 'react-router-dom';
import './style/HomePage.css';

function HomePage() {
    const [gameData, setGameData] = useState(null);
    const [answer, setAnswer] = useState(null);
    const [gameSessionLoading, setGameSessionLoading] = useState(false);
    const [gameError, setGameError] = useState(null);
    const [gameOver, setGameOver] = useState(false);
    const [gameScore, setGameScore] = useState(0);
    const [timeLeft, setTimeLeft] = useState(0);
    const [intervalId, setIntervalId] = useState(null);
    const [formattedTime, setFormattedTime] = useState("");
    const [showScoreModal, setShowScoreModal] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (!token) {
            navigate('/login');  
        }
    }, [navigate]);

    const formatTime = (timeInNano) => {
        const timeInSeconds = timeInNano / 1000000000;
        const minutes = Math.floor(timeInSeconds / 60);
        const seconds = Math.floor(timeInSeconds % 60);
        return `${minutes}:${seconds < 10 ? '0' + seconds : seconds}`;
    };

    const handleCreateGameSession = async () => {
        setGameSessionLoading(true);
        setGameError(null);
        try {
            const data = await createGameSession(4);  
            setGameData(data);
            setFormattedTime(formatTime(data.time_left));
            startTimer(data.time_left);  
        } catch (err) {
            setGameError('Failed to create game session');
        } finally {
            setGameSessionLoading(false);
        }
    };

    const startTimer = (timeInNano) => {
        const timeInSeconds = timeInNano / 1000000000;
        setTimeLeft(timeInSeconds);

        const id = setInterval(() => {
            setTimeLeft((prevTime) => {
                if (prevTime <= 0) {
                    clearInterval(id);
                    handleFinishGame();
                    return 0;
                }
                return prevTime - 1;
            });
        }, 1000);
        setIntervalId(id);
    };

    const handleSubmitAnswer = async () => {
        if (!answer) return;
        try {
            const data = await submitAnswer(gameData.session_id, 4, answer);
            setGameData(data);
            setGameScore(data.score);
            setAnswer(null);
        } catch (err) {
            setGameError('Failed to submit answer');
        }
    };

    const handleFinishGame = async () => {
        try {
            const result = await finishGame(gameData.session_id, 4);
            setGameOver(true);
            clearInterval(intervalId); 
            setShowScoreModal(true);  
        } catch (err) {
            setGameError('Failed to finish game');
            setTimeout(() => {
                navigate(0);  
            }, 2000);
        }
    };

    useEffect(() => {
        if (timeLeft > 0) {
            setFormattedTime(formatTime(timeLeft * 1000000000));
        }
    }, [timeLeft]);

    const handleRestartGame = () => {
        setGameData(null);
        setAnswer(null);
        setGameScore(0);
        setTimeLeft(0);
        setFormattedTime("");
        setShowScoreModal(false);
        setGameOver(false);
        handleCreateGameSession(); 
    };

    const handleExitGame = () => {
        localStorage.removeItem('authToken');
        navigate('/login');
    };

    return (
        <div className="home-container">
            <h1 className="home-title">Welcome to the MathArena</h1>

            {/* Centered Buttons for Creating Game and Logging Out */}
            {!gameData && !gameOver && (
                <div className="center-buttons">
                    <button
                        className="create-game-button"
                        onClick={handleCreateGameSession}
                        disabled={gameSessionLoading}
                    >
                        {gameSessionLoading ? 'Creating Game...' : 'Create Your Math Game'}
                    </button>

                    <button
                        className="logout-button"
                        onClick={handleExitGame}
                    >
                        Log out
                    </button>
                </div>
            )}

            {/* Game Popup */}
            {gameData && !gameOver && (
                <div className="game-popup-overlay">
                    <div className="game-popup">
                        <h2>Math Expression: {gameData.expression}</h2>
                        <p>Time Left: {formattedTime}</p> 
                        <p>Score: {gameScore}</p>

                        <input className="input_answer"
                            type="number"
                            value={answer || ''}
                            onChange={(e) => setAnswer(Number(e.target.value))}
                            placeholder="Enter your answer"
                        />
                        <button className="submit-answer-btn" onClick={handleSubmitAnswer}>Submit Answer</button>
                        <button className="finish-game-btn" onClick={handleFinishGame}>Finish Game</button>
                    </div>
                </div>
            )}

            {gameError && <p className="error-message">{gameError}</p>}

            {/* Score Modal that appears when the game is over */}
            {showScoreModal && (
                <div className="score-modal">
                    <h2>Game Over!</h2>
                    <p>Your final score is: {gameScore}</p>
                    <button onClick={handleRestartGame}>Play Again</button>
                    <button onClick={handleExitGame}>Log Out</button>
                </div>
            )}

            {/* Logout button that appears after the game ends */}
            {gameOver && !showScoreModal && (
                <button
                    className="logout-button"
                    onClick={handleExitGame}
                >
                    Log out
                </button>
            )}
        </div>
    );
}

export default HomePage;

import { useState, useEffect, useCallback } from 'react';
import { LoginButton, TelegramAuthData } from '@telegram-auth/react';

function Auth() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [showLoginButton, setShowLoginButton] = useState(false);
    const [key, setKey] = useState(0);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
        } else {
            setShowLoginButton(true);
        }
    }, []);

    const handleAuthCallback = useCallback(async (data: TelegramAuthData) => {
        try {
            const response = await fetch('http://localhost:4000/auth/telegram/callback', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const result = await response.json();

            if (result.success) {
                localStorage.setItem('token', result.token);
                setIsAuthenticated(true);
                setShowLoginButton(false);
                window.location.reload()
            } else {
                console.error('Authentication failed:', result.message);
            }
        } catch (error) {
            console.error('Error during authentication:', error);
        }
    }, []);

    const handleLogout = useCallback(() => {
        localStorage.removeItem('token');
        setIsAuthenticated(false);
        setShowLoginButton(true);
        setKey(prevKey => prevKey + 1);
        window.location.reload()
    }, []);

    if (isAuthenticated) {
        return (
            <a className="header__auth-logout" onClick={handleLogout}>Выйти</a>
        );
    }

    return (
        <div key={key}>
            {showLoginButton && (
                <LoginButton
                    botUsername="kinogolang_bot"
                    buttonSize="medium"
                    cornerRadius={10}
                    showAvatar={false}
                    lang="ru"
                    onAuthCallback={handleAuthCallback}
                />
            )}
        </div>
    );
}

export default Auth;
import { useState, useEffect, useCallback } from 'react';
import { LoginButton, TelegramAuthData } from '@telegram-auth/react';
import Cookies from 'js-cookie';
import { GetIsAdmin } from '@components/JwtDecode.tsx';

function Auth() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userAdmin, setUserAdmin] = useState(false);
    const isAdmin = GetIsAdmin();

    useEffect(() => {
        const token = Cookies.get('token');
        setIsAuthenticated(!!token);
    }, []);

    useEffect(() => {
        if (isAdmin !== undefined) {
            setUserAdmin(isAdmin);
        }
    }, [isAdmin]);

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
                console.log(`HTTP error! Status: ${response.status}`);
            }

            const result = await response.json();

            if (result.success) {
                Cookies.set('token', result.token);
                setIsAuthenticated(true);
                window.location.reload();
            } else {
                console.error('Authentication failed:', result.message);
            }
        } catch (error) {
            console.error('Error during authentication:', error);
        }
    }, []);

    const handleLogout = useCallback(() => {
        Cookies.remove('token');
        setIsAuthenticated(false);
        window.location.reload();
    }, []);

    return (
        <>
            {isAuthenticated ? (
                <>
                    {userAdmin && (
                        <a className="header__auth-admin" href="/studio">Творческая студия</a>
                    )}
                    <a className="header__auth-logout" onClick={handleLogout}>Выйти</a>
                </>
            ) : (
                <div>
                    <LoginButton
                        botUsername="kinogolang_bot"
                        buttonSize="medium"
                        cornerRadius={10}
                        showAvatar={false}
                        lang="ru"
                        onAuthCallback={handleAuthCallback}
                    />
                </div>
            )}
        </>
    );
}

export default Auth;
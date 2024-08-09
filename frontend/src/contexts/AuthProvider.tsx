import { useState, useEffect, ReactNode } from 'react';
import Cookies from 'js-cookie';
import {jwtDecode} from "jwt-decode";
import { AuthContext, JwtPayload } from './AuthContext';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userAdmin, setUserAdmin] = useState(false);
    const [userId, setUserId] = useState<number | null>(null);

    useEffect(() => {
        const token = Cookies.get('token');
        if (token) {
            setIsAuthenticated(true);
            try {
                const payload = jwtDecode<JwtPayload>(token);
                setUserAdmin(payload.isAdmin);
                setUserId(payload.id);
            } catch (error) {
                console.error('Invalid token', error);
                setIsAuthenticated(false);
                setUserAdmin(false);
                setUserId(null);
            }
        }
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, userAdmin, userId, setIsAuthenticated, setUserAdmin, setUserId }}>
            {children}
        </AuthContext.Provider>
    );
};
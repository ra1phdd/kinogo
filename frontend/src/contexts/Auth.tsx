import { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import Cookies from 'js-cookie';
import {jwtDecode} from "jwt-decode";

interface JwtPayload {
    id: number;
    isAdmin: boolean;
}

interface AuthContextType {
    isAuthenticated: boolean;
    userAdmin: boolean;
    userId: number | null;
    setIsAuthenticated: (auth: boolean) => void;
    setUserAdmin: (isAdmin: boolean) => void;
    setUserId: (id: number | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

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

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
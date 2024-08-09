import { createContext } from 'react';

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

export {AuthContext};
export type { AuthContextType, JwtPayload };


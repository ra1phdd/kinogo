import Cookies from "js-cookie";
import {jwtDecode} from "jwt-decode";
import {useEffect, useState} from "react";

interface JwtPayload {
    id: number;
    isAdmin: boolean;
}

export function GetUserId(): number | undefined {
    const [userId, setUserId] = useState<number>();

    useEffect(() => {
        const token = Cookies.get('token');
        if (token) {
            try {
                const payload = jwtDecode<JwtPayload>(token);
                setUserId(payload.id);
            } catch (error) {
                console.error('Invalid token', error);
                setUserId(0);
            }
        }
    }, []);

    return userId;
}

export function GetIsAdmin(): boolean | undefined {
    const [isAdmin, setIsAdmin] = useState<boolean>();

    useEffect(() => {
        const token = Cookies.get('token');
        if (token) {
            try {
                const payload = jwtDecode<JwtPayload>(token);
                setIsAdmin(payload.isAdmin);
            } catch (error) {
                console.error('Invalid token', error);
                setIsAdmin(false);
            }
        }
    }, []);

    return isAdmin;
}
import React, { useEffect } from 'react';
import {metricAvgTimeOnSite} from "@components/gRPC.tsx";

const TimeTracker: React.FC = () => {
    const getCurrentTime = (): number => new Date().getTime();

    const setLastTime = (): void => {
        localStorage.setItem('lastTime', getCurrentTime().toString());
    };

    const getLastTime = (): number => {
        return parseInt(localStorage.getItem('lastTime') || '0');
    };

    const setFirstTime = (): void => {
        if (!localStorage.getItem('firstTime')){
            localStorage.setItem('firstTime', getCurrentTime().toString());
        }
    };

    const getFirstTime = (): number => {
        return parseInt(localStorage.getItem('firstTime') || '0');
    };

    const sendDataToBackend = (timeSpent: number): void => {
        try {
            metricAvgTimeOnSite(timeSpent);
        } catch (error) {
            console.error('Error sending time spent:', error);
        }
    };

    const checkAndSendData = (): void => {
        const firstTime = getFirstTime();
        const currentTime = getCurrentTime();
        const lastTime = getLastTime();
        const timeSpent = currentTime - lastTime;
        const timestamp = lastTime - firstTime;

        if ((timeSpent >= 30000) && (localStorage.getItem('lastTime') && (localStorage.getItem('lastTime')))) {
            sendDataToBackend(timestamp);
            localStorage.removeItem('lastTime');
            localStorage.removeItem('lastTime');
        }
    };

    // useEffect для установки времени входа и проверки перед выходом
    useEffect(() => {
        setFirstTime();
        checkAndSendData();
    }, []);

    useEffect(() => {
        const intervalId = setInterval(() => {
            setLastTime();
        }, 5000);

        return () => {
            clearInterval(intervalId);
        };
    }, []);

    return <></>;
};

export default TimeTracker;

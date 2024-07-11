import React, { useEffect } from 'react';
import {metricAvgTimeOnSite} from "@components/gRPC.tsx";

const TimeTracker: React.FC = () => {
    const getCurrentTime = (): number => new Date().getTime();
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
        const timeSpent = currentTime - firstTime;

        sendDataToBackend(timeSpent);
        localStorage.removeItem('firstTime');
    };

    useEffect(() => {
        setFirstTime();

        window.addEventListener('beforeunload', checkAndSendData);

        return () => {
            window.removeEventListener('beforeunload', checkAndSendData);
        };
    }, []);


    return <></>;
};

export default TimeTracker;

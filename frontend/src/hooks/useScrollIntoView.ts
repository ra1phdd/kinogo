import { useEffect, useRef } from 'react';

export const useScrollIntoView = (condition: boolean) => {
    const formRef = useRef<HTMLFormElement>(null);

    useEffect(() => {
        if (condition && formRef.current) {
            formRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }, [condition]);

    return formRef;
};

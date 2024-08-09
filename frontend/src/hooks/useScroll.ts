import { useCallback, useEffect } from "react";

const useScroll = (loadMovies: () => void) => {
    const handleScroll = useCallback(() => {
        if (window.innerHeight + document.documentElement.scrollTop >=
            document.documentElement.offsetHeight - 200) {
            loadMovies();
        }
    }, [loadMovies]);

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [handleScroll]);
};

export default useScroll;
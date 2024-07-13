import { useState, useEffect } from 'react';
import { BestMovie, getBestMovie } from '@components/gRPC.tsx';

const useBestMovie = () => {
    const [movie, setMovie] = useState<BestMovie | null>(null);

    useEffect(() => {
        const fetchMovie = async () => {
            try {
                const movieResponse = await getBestMovie();
                setMovie(movieResponse[0]);
            } catch (error) {
                console.error('Ошибка при получении фильма:', error);
            }
        };

        fetchMovie();
    }, []);

    return movie;
};

export default useBestMovie;

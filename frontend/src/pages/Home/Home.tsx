import React, { useCallback, useEffect, useState } from "react";
import { Movies, getMovies } from '@components/gRPC.tsx';
import MovieCard from "@components/common/MovieCard.tsx";
import {RpcError} from "grpc-web";
import '@assets/styles/pages/home.css'

// Компонент Home
const Home: React.FC = () => {
    const [movies, setMovies] = useState<Movies[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(1);

    const loadMovies = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const newMovies = await getMovies(15, page);
                setMovies((prevMovies) => [...prevMovies, ...newMovies]);
                setPage((prevPage) => prevPage + 1);
            } catch (error) {
                if (error instanceof RpcError) {
                    if (error.code === 5) {
                        setEnd(true);
                    } else {
                        console.error('Error fetching more movies:', error);
                    }
                } else {
                    console.error('Unexpected error:', error);
                }
            } finally {
                setLoading(false);
            }
        }
    }, [loading, end]);

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

    useEffect(() => {
        loadMovies();
    }, []);

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    {movies.map((movie, index) => (
                        <MovieCard key={movie.id} movie={movie} index={index} limit={15}  />
                    ))}
                </div>
            </div>
            {loading && <div>Загрузка...</div>}
        </>
    );
};

export default Home;

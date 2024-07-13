import React, { useCallback, useEffect, useState } from "react";
import SearchAside from "@components/specific/Search.tsx";
import FilterAside from "@components/specific/Filter.tsx";
import BestMovieAside from "@components/specific/BestMovie.tsx";
import { getMoviesByTypeMovie } from '@components/gRPC.tsx';
import { useLocation } from "react-router-dom";
import MovieCard from "@components/common/MovieCard.tsx";
import {RpcError} from "grpc-web";
import '@assets/styles/pages/movies.css';
import type { Movies } from '@components/gRPC.tsx';

// Компонент Home
const Movies: React.FC = () => {
    const [movies, setMovies] = useState<Movies[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(1);
    const location = useLocation();

    const fetchFunctions: Record<string, () => Promise<Movies[]>> = {
        '/films': () => getMoviesByTypeMovie(12, page, 1),
        '/cartoons': () => getMoviesByTypeMovie(12, page, 2),
        '/telecasts': () => getMoviesByTypeMovie(12, page, 3),
    };
    const fetchFunction = fetchFunctions[location.pathname] || fetchFunctions['/films'];

    const loadMovies = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const newMovies = await fetchFunction();
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
    }, [loading, end, fetchFunction]);


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
                        <MovieCard key={movie.id} movie={movie} index={index} limit={12} />
                    ))}
                    {loading && <p>Загрузка...</p>}
                    {movies.length == 0 && <p>Данных в БД нет</p>}
                </div>
            </div>
            <aside>
                <SearchAside />
                <FilterAside />
                <BestMovieAside />
            </aside>
        </>
    );
};

export default Movies;

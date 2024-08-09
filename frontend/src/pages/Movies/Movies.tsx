import React, { useEffect } from "react";
import SearchAside from "@components/specific/Search.tsx";
import FilterAside from "@components/specific/Filter.tsx";
import BestMovieAside from "@components/specific/BestMovie.tsx";
import {useLocation} from "react-router-dom";
import MovieCard from "@components/common/MovieCard.tsx";
import '@assets/styles/pages/movies.css';
import useScroll from "@/hooks/useScroll.ts";
import useTypeMovies from "@/hooks/fetchMovies/useTypeMovies.ts";

// Компонент Home
const Movies: React.FC = () => {
    const location = useLocation();

    const { movies, loading, loadMovies } = useTypeMovies(1, location.pathname);

    useScroll(loadMovies);

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

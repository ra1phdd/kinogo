import React, { useEffect, useState } from "react";
import { BestMovie, getBestMovie } from "@components/gRPC.tsx";
import AnimateElement from "@components/AnimateElement.tsx";

// Пропсы для компонента MovieCard
interface MovieCardProps {
    movie: BestMovie;
}

// Компонент MovieCard
const MovieCard: React.FC<MovieCardProps> = ({ movie }) => {
    useEffect(() => {
        const asideBestMovie = document.querySelector<HTMLElement>(".aside__bestmovie");

        if (asideBestMovie !== null) {
            AnimateElement(asideBestMovie, "animate__fadeInRight", 450);
        }
    }, []);

    return (
        <>
            <h3>Популярный фильм</h3>
            <div className="bestmovie__item">
                <img src={movie.poster} alt=""/>
                <p>{movie.title} ({movie.releaseDate})</p>
                <button>
                    <div id="circle"></div>
                    <a href={`/id/${movie.id}`}>Смотреть</a>
                </button>
            </div>
        </>
    );
};

// Компонент BestMovieAside
const BestMovieAside: React.FC = () => {
    const [movies, setMovies] = useState<BestMovie[]>([]);

    useEffect(() => {
        const fetchMovies = async () => {
            try {
                const moviesResponse = await getBestMovie();
                setMovies(moviesResponse);
            } catch (error) {
                console.error('Ошибка при получении фильмов:', error);
            }
        };

        fetchMovies();
    }, []);

    if (movies.length === 0) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="aside__bestmovie">
            {movies.map((movie: BestMovie) => (
                <MovieCard key={movie.id} movie={movie} />
            ))}
        </div>
    );
}

export default BestMovieAside;

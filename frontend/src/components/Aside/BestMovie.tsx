import React, { useEffect, useState } from "react";
import { BestMovie, getBestMovie } from "@components/gRPC.tsx";
import AnimateElement from "@components/AnimateElement.tsx";

// Компонент BestMovieAside
const BestMovieAside: React.FC = () => {
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

    useEffect(() => {
        const asideBestMovie = document.querySelector<HTMLElement>(".aside__bestmovie");

        if (asideBestMovie !== null) {
            AnimateElement(asideBestMovie, "animate__fadeInRight", 150);
        }
    }, [movie]);

    if (!movie) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="aside__bestmovie">
            <h3>Популярный фильм</h3>
            <div className="bestmovie__item">
                <img src={movie.poster} alt=""/>
                <p>{movie.title} ({movie.releaseDate})</p>
                <button>
                    <div id="circle"></div>
                    <a href={`/id/${movie.id}`}>Смотреть</a>
                </button>
            </div>
        </div>
    );
}

export default BestMovieAside;
import React, { useEffect } from "react";
import AnimateElement from "@components/common/AnimateElement.tsx";
import useBestMovie from "@/hooks/useBestMovie.ts";

// Компонент BestMovieAside
const BestMovieAside: React.FC = () => {
    const movie = useBestMovie();

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
import React, {useCallback, useEffect, useRef} from "react";
import AnimateElement from "@components/common/AnimateElement.tsx";
import {Movies} from "@components/gRPC.tsx";

interface MovieCardProps {
    movie: Movies;
    index: number;
    limit: number
}

const MovieCard: React.FC<MovieCardProps> = React.memo(({ movie, index, limit }) => {
    const cardRef = useRef<HTMLDivElement>(null);

    const animate = useCallback(() => {
        let delay: number;
        if (cardRef.current && !cardRef.current.classList.contains('animate__animated')) {
            const adjustedIndex = index % limit;
            delay = adjustedIndex * 75;
            AnimateElement(cardRef.current, "animate__fadeInUp", delay);
        }
    }, [index]);


    useEffect(() => {
        animate();
    }, [animate]);

    return (
        <div ref={cardRef} className={`card-${movie.id} card`}>
            <a href={`/id/${movie.id}`} data-tilt="">
                <div className="card__poster">
                    <img src={movie.poster} alt="Постер" />
                </div>
                <div className="card__info">
                    <h2 className="card__info-title">{movie.title} ({movie.releaseDate})</h2>
                    {movie.scoreKP > 7 && (
                        <div className="card__info-rating" style={{backgroundColor: "green"}}>
                            <p>{movie.scoreKP}</p>
                        </div>
                    )}
                    {movie.scoreKP > 4 && movie.scoreKP < 7 && (
                        <div className="card__info-rating" style={{backgroundColor: "gray"}}>
                            <p>{movie.scoreKP}</p>
                        </div>
                    )}
                    {movie.scoreKP < 4 && (
                        <div className="card__info-rating" style={{backgroundColor: "red"}}>
                            <p>{movie.scoreKP}</p>
                        </div>
                    )}
                    <div className="card__info-genres">
                        <p>{movie.genres}</p>
                    </div>
                </div>
            </a>
        </div>
    );
});

export default MovieCard
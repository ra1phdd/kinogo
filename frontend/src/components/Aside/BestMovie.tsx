// Интерфейс для данных фильма
import React, {useEffect, useRef, useState} from "react";

interface Movie {
    Id: number;
    Title: string;
    Description: string;
    Country: string;
    ReleaseDate: number;
    TimeMovie: number;
    ScoreKP: number;
    ScoreIMDB: number;
    Poster: string;
    TypeMovie: string;
    Views: number;
    Likes: number;
    Dislikes: number;
    Genres: string;
}

// Пропсы для компонента MovieCard
interface MovieCardProps {
    movie: Movie;
}

// Компонент MovieCard
const MovieCard: React.FC<MovieCardProps> = ({ movie }) => {
    const cardRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (cardRef.current) {
            cardRef.current.classList.add("animate__animated", "animate__fadeInLeft", "animate__faster");
        }
    }, []);

    return (
        <div ref={cardRef} className="aside__bestmovie">
            <h3>Популярный фильм</h3>
            <div className="bestmovie__item">
                <img src={movie.Poster} alt=""/>
                <p>{movie.Title} ({movie.ReleaseDate})</p>
                <button>
                    <div id="circle"></div>
                    <a href={`/id/${movie.Id}`}>Смотреть</a>
                </button>
            </div>
        </div>
    );
};

const BestMovieAside: React.FC = () =>  {
    const [movie, setMovies] = useState<Movie | null>(null);

    useEffect(() => {
        fetch('http://localhost:4000/api/v1/content/best')
            .then(response => response.json())
            .then(data => setMovies(data))
            .catch(error => console.error('Ошибка загрузки данных:', error));
    }, []);

    if (!movie) {
        return <div>Загрузка...</div>;
    }

    return (
        <MovieCard key={movie.Id} movie={movie} />
    )
}

export default BestMovieAside

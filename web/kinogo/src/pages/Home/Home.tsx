import React, {useEffect, useRef, useState} from "react";
import AnimateElement from "../../components/AnimateElement.tsx";
import SearchAside from "../../components/Aside/Search.tsx";
import FilterAside from "../../components/Aside/Filter.tsx";
import BestMovieAside from "../../components/Aside/BestMovie.tsx";

// Интерфейс для данных фильма
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
        <div ref={cardRef} className={`card-${movie.Id} card`}>
            <a href={`/id/${movie.Id}`} data-tilt="">
                <div className="card__poster">
                    <img src={movie.Poster} alt="Постер" />
                </div>
                <div className="card__info">
                    <h2 className="card__info-title">{movie.Title} ({movie.ReleaseDate})</h2>
                    <p className="card__info-description">{movie.Description}</p>
                    <div className="card__info-ratings">
                        <div className="rating__kinopoisk rating">
                            <p>Кинопоиск {movie.ScoreKP}</p>
                        </div>
                        <div className="rating__imdb rating">
                            <p>IMDb {movie.ScoreIMDB}</p>
                        </div>
                    </div>
                    <div className="card__info-genres">
                        <p>{movie.Genres}</p>
                    </div>
                </div>
            </a>
        </div>
    );
};

// Компонент Home
const Home: React.FC = () => {
    const [movies, setMovies] = useState<Movie[]>([]);

    // Получение фильмов по API
    useEffect(() => {
        fetch('http://localhost:4000/api/v1/contents')
            .then(response => response.json())
            .then(data => setMovies(data))
            .catch(error => console.error('Ошибка загрузки данных:', error));
    }, []);

    if (!movies) {
        return <div>Загрузка...</div>;
    }

    // Анимация
    useEffect(() => {
        const cards = document.querySelectorAll<HTMLElement>(".card");
        let delay = 0;

        cards.forEach((card) => {
            AnimateElement(card, "animate__fadeInLeft", delay);
            delay += 150;
        });

        const asideSearch = document.querySelector<HTMLElement>(".aside__search")
        const asideFilter = document.querySelector<HTMLElement>(".aside__filter")
        const asideBestMovie = document.querySelector<HTMLElement>(".aside__bestmovie")

        if (asideSearch != null){
            AnimateElement(asideSearch, "animate__fadeInRight", 0);
        }
        if (asideFilter != null){
            AnimateElement(asideFilter, "animate__fadeInRight", 150);
        }
        if (asideBestMovie != null){
            AnimateElement(asideBestMovie, "animate__fadeInRight", 300);
        }
    }, [movies]);

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    {movies.map(movie => (
                        <MovieCard key={movie.Id} movie={movie} />
                    ))}
                </div>
            </div>
            <aside>
                <SearchAside/>
                <FilterAside/>
                <BestMovieAside/>
            </aside>
        </>
    );
};

export default Home;

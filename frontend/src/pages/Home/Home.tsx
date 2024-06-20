import React, {useEffect, useRef, useState} from "react";
import AnimateElement from "@components/AnimateElement.tsx";
import SearchAside from "@components/Aside/Search.tsx";
import FilterAside from "@components/Aside/Filter.tsx";
import BestMovieAside from "@components/Aside/BestMovie.tsx";
import { getMovies } from '@components/gRPC.tsx';

// Интерфейс для данных фильма
interface Movie {
    id: number;
    title: string;
    description: string;
    releaseDate: number;
    scoreKP: number;
    scoreIMDB: number;
    poster: string;
    typeMovie: number;
    genres: string;
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
        <div ref={cardRef} className={`card-${movie.id} card`}>
            <a href={`/id/${movie.id}`} data-tilt="">
                <div className="card__poster">
                    <img src={movie.poster} alt="Постер" />
                </div>
                <div className="card__info">
                    <h2 className="card__info-title">{movie.title} ({movie.releaseDate})</h2>
                    <p className="card__info-description">{movie.description}</p>
                    <div className="card__info-ratings">
                        <div className="rating__kinopoisk rating">
                            <p>Кинопоиск {movie.scoreKP}</p>
                        </div>
                        <div className="rating__imdb rating">
                            <p>IMDb {movie.scoreIMDB}</p>
                        </div>
                    </div>
                    <div className="card__info-genres">
                        <p>{movie.genres}</p>
                    </div>
                </div>
            </a>
        </div>
    );
};

// Компонент Home
const Home: React.FC = () => {
    const [movies, setMovies] = useState<Movie[]>([]);

    useEffect(() => {
        const fetchMovies = async () => {
            try {
                const moviesResponse = await getMovies(10, 1); // Fetch first 10 movies
                setMovies(moviesResponse);
            } catch (error) {
                console.error('Error fetching movies:', error);
            }
        };

        fetchMovies();
    }, []);

    useEffect(() => {
        // Function to animate elements
        const animateElements = () => {
            const cards = document.querySelectorAll<HTMLElement>(".card");
            let delay = 0;

            cards.forEach((card) => {
                AnimateElement(card, "animate__fadeInLeft", delay);
                delay += 150;
            });

            const asideSearch = document.querySelector<HTMLElement>(".aside__search");
            const asideFilter = document.querySelector<HTMLElement>(".aside__filter");
            const asideBestMovie = document.querySelector<HTMLElement>(".aside__bestmovie");

            if (asideSearch !== null) {
                AnimateElement(asideSearch, "animate__fadeInRight", 0);
            }
            if (asideFilter !== null) {
                AnimateElement(asideFilter, "animate__fadeInRight", 150);
            }
            if (asideBestMovie !== null) {
                AnimateElement(asideBestMovie, "animate__fadeInRight", 300);
            }
        };

        animateElements(); // Initial animation on component mount

        // Re-run animation whenever movies change
        if (movies.length > 0) {
            animateElements();
        }
    }, [movies]);

    if (movies.length === 0) {
        return <div>Загрузка...</div>;
    }

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    {movies.map(movie => (
                        <MovieCard key={movie.id} movie={movie} />
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

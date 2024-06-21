import React, {useCallback, useEffect, useRef, useState} from "react";
import AnimateElement from "@components/AnimateElement.tsx";
import SearchAside from "@components/Aside/Search.tsx";
import FilterAside from "@components/Aside/Filter.tsx";
import BestMovieAside from "@components/Aside/BestMovie.tsx";
import { Movies, getMovies } from '@components/gRPC.tsx';

// Пропсы для компонента MovieCard
interface MovieCardProps {
    movie: Movies;
    index: number;
}

// Компонент MovieCard
const MovieCard: React.FC<MovieCardProps> = React.memo(({ movie, index }) => {
    const cardRef = useRef<HTMLDivElement>(null);

    const animate = useCallback(() => {
        if (cardRef.current) {
            AnimateElement(cardRef.current, "animate__fadeInLeft", index * 150);
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
});

// Компонент Home
const Home: React.FC = () => {
    const [movies, setMovies] = useState<Movies[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    const fetchMovies = useCallback(async () => {
        try {
            const moviesResponse = await getMovies(10, 1); // Fetch first 10 movies
            setMovies(moviesResponse);
        } catch (error) {
            console.error('Error fetching movies:', error);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchMovies();
    }, [fetchMovies]);

    if (loading) {
        return <div>Загрузка...</div>;
    }

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    {movies.map((movie, index) => (
                        <MovieCard key={movie.id} movie={movie} index={index}/>
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
}

export default Home;
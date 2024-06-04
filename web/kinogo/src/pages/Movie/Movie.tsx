import React, {useEffect, useRef, useState} from 'react';
import { useParams } from 'react-router-dom';
import AnimateElement from '../../components/AnimateElement';
import SearchAside from '../../components/Aside/Search';
import FilterAside from '../../components/Aside/Filter';
import BestMovieAside from '../../components/Aside/BestMovie';
import '../../assets/css/src/movie.css'
import VideoPlayer from "../../components/HLSPlayer.tsx";

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
        <>
            <div ref={cardRef} className="section__movie">
                <div className="section__movie-img">
                    <img src={movie.Poster} alt="Постер фильма"/>
                </div>
                <div className="section__movie-info">
                    <h1>{movie.Title}</h1>
                    <p>{movie.Description}</p>
                    <div className="section__movie-main">
                        <div className="section__movie-buttons">
                            <div className="button__like">
                                <form action="/like" method="post">
                                    <input type="hidden" name="id" value={movie.Id} className="likeValue"/>
                                    <button type="submit" className="likeButton">Поставить лайк</button>
                                </form>
                            </div>
                            <div className="button__dislike">
                                <form action="/dislike" method="post">
                                    <input type="hidden" name="id" value={movie.Id} className="dislikeValue"/>
                                    <button type="submit" className="dislikeButton">Поставить дизлайк</button>
                                </form>
                            </div>
                        </div>
                        <div className="section__movie-ratings">
                            <div className="rating__kinopoisk rating">
                                <p>Кинопоиск {movie.ScoreKP}</p>
                            </div>
                            <div className="rating__imdb rating">
                                <p>IMDb {movie.ScoreIMDB}</p>
                            </div>
                        </div>
                    </div>
                    <h3>О фильме</h3>
                    <table>
                        <tbody>
                            <tr>
                                <td className="titleInfo">Год выхода</td>
                                <td className="itemInfo">{movie.ReleaseDate}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Страна</td>
                                <td className="itemInfo">{movie.Country}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Жанр</td>
                                <td className="itemInfo">{movie.Genres}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Длительность</td>
                                <td className="itemInfo">{movie.TimeMovie} минут</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <div className="section__player">
                <VideoPlayer/>
            </div>
        </>
    )
        ;
};

// Компонент Movie
const Movie: React.FC = () => {
    let {id} = useParams<{ id: string }>();
    const [movie, setMovie] = useState<Movie | null>(null);

    useEffect(() => {
        fetch(`http://localhost:4000/api/v1/movies/${id}`)
            .then(response => response.json())
            .then(data => setMovie(data))
            .catch(error => console.error('Ошибка загрузки данных:', error));
    }, []);

    // Анимация
    useEffect(() => {
        const sectionMovie = document.querySelector<HTMLElement>(".section__movie")
        const sectionPlayer = document.querySelector<HTMLElement>(".section__player")
        const asideSearch = document.querySelector<HTMLElement>(".aside__search")
        const asideFilter = document.querySelector<HTMLElement>(".aside__filter")
        const asideBestMovie = document.querySelector<HTMLElement>(".aside__bestmovie")

        if (sectionMovie) {
            AnimateElement(sectionMovie, "animate__fadeInLeft", 0);
        }
        if (sectionPlayer) {
            AnimateElement(sectionPlayer, "animate__fadeInLeft", 150);
        }
        if (asideSearch) {
            AnimateElement(asideSearch, "animate__fadeInRight", 0);
        }
        if (asideFilter) {
            AnimateElement(asideFilter, "animate__fadeInRight", 150);
        }
        if (asideBestMovie) {
            AnimateElement(asideBestMovie, "animate__fadeInRight", 300);
        }
    }, [movie]);

    if (!movie) {
        return <div>Загрузка...</div>;
    }

    return (
        <>
            <div className="sections">
                <div className="section__videos">
                    <MovieCard key={movie.Id} movie={movie} />
                </div>
            </div>
            <aside>
                <SearchAside />
                <FilterAside />
                <BestMovieAside/>
            </aside>
        </>
    );
};

export default Movie;

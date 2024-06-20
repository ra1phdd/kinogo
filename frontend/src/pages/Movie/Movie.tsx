import React, {useEffect, useRef, useState} from 'react';
import { useParams } from 'react-router-dom';
import AnimateElement from '@components/AnimateElement';
import SearchAside from '@components/Aside/Search';
import FilterAside from '@components/Aside/Filter';
import BestMovieAside from '@components/Aside/BestMovie';
import '@assets/css/src/movie.css'
import VideoPlayer from "@components/HLSPlayer.tsx";
import { Movie, getMovieById } from "@components/gRPC.tsx";

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
                    <img src={movie.poster} alt="Постер фильма"/>
                </div>
                <div className="section__movie-info">
                    <h1>{movie.title}</h1>
                    <p>{movie.description}</p>
                    <div className="section__movie-main">
                        <div className="section__movie-buttons">
                            <div className="button__like">
                                <form action="/like" method="post">
                                    <input type="hidden" name="id" value={movie.id} className="likeValue"/>
                                    <button type="submit" className="likeButton">Поставить лайк</button>
                                </form>
                            </div>
                            <div className="button__dislike">
                                <form action="/dislike" method="post">
                                    <input type="hidden" name="id" value={movie.id} className="dislikeValue"/>
                                    <button type="submit" className="dislikeButton">Поставить дизлайк</button>
                                </form>
                            </div>
                        </div>
                        <div className="section__movie-ratings">
                            <div className="rating__kinopoisk rating">
                                <p>Кинопоиск {movie.scoreKP}</p>
                            </div>
                            <div className="rating__imdb rating">
                                <p>IMDb {movie.scoreIMDB}</p>
                            </div>
                        </div>
                    </div>
                    <h3>О фильме</h3>
                    <table>
                        <tbody>
                            <tr>
                                <td className="titleInfo">Год выхода</td>
                                <td className="itemInfo">{movie.releaseDate}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Страна</td>
                                <td className="itemInfo">{movie.country}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Жанр</td>
                                <td className="itemInfo">{movie.genres}</td>
                            </tr>
                            <tr>
                                <td className="titleInfo">Длительность</td>
                                <td className="itemInfo">{movie.timeMovie} минут</td>
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
const GetMovie: React.FC = () => {
    const { id } = useParams<{ id: string }>(); // Изменяем тип id на string
    const idNumber = Number(id);
    const [movie, setMovie] = useState<Movie | null>(null);

    useEffect(() => {
        const fetchMovieById = async () => {
            try {
                const moviesResponse = await getMovieById(idNumber);
                setMovie(moviesResponse);
            } catch (error) {
                console.error('Error fetching movie:', error);
            }
        };

        fetchMovieById();
    }, [idNumber]);

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
                    <MovieCard key={movie.id} movie={movie} />
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

export default GetMovie;

import React, {useCallback, useEffect, useRef, useState} from 'react';
import { useParams } from 'react-router-dom';
import AnimateElement from '@components/common/AnimateElement.tsx';
import '@assets/styles/pages/movie.css';
import VideoPlayer from "@components/specific/HLSPlayer.tsx";
import { Movie, getMovieById } from "@components/gRPC.tsx";
import CommentsComponent from "@components/specific/Comments/CommentsComponent.tsx";
import Cookies from "js-cookie";

// Пропсы для компонента MovieCard
interface MovieCardProps {
    movie: Movie;
}

// Компонент MovieCard
const MovieDetails: React.FC<MovieCardProps> = ({ movie }) => {
    const cardRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (cardRef.current) {
            cardRef.current.classList.add("animate__animated", "animate__fadeInLeft", "animate__faster");
        }
    }, []);

    return (
        <div ref={cardRef} className="section__movie">
            <div className="section__movie-img">
                <img src={movie.poster} alt="Постер фильма"/>
            </div>
            <div className="section__movie-info">
                <h1>{movie.title}</h1>
                <p>{movie.description}</p>
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
    );
};

const MovieCard: React.FC<MovieCardProps> = ({ movie }) => {
    const [hasToken, setHasToken] = useState(false);

    useEffect(() => {
        const token = Cookies.get('token');
        setHasToken(!!token);
    }, []);

    return (
        <>
            <MovieDetails movie={movie} />
            {hasToken ? (
                <div className="section__player">
                    <VideoPlayer id={movie.id} title={movie.title}/>
                </div>
            ) : (
                <div className="section__player out">
                    {!hasToken && (
                        <div className="section__player-hide">Для просмотра этого контента авторизуйтесь на сайте</div>
                    )}
                    <VideoPlayer id={movie.id} title={movie.title}/>
                </div>
            )}
            <div className="section__comments">
                <h2>Комментарии</h2>
                <CommentsComponent id={movie.id}/>
            </div>
        </>
    );
};

// Компонент Movie
const GetMovie: React.FC = () => {
    const {id} = useParams<{ id: string }>();
    const [movie, setMovie] = useState<Movie | null>(null);

    const fetchMovieById = useCallback(async () => {
        try {
            const movieId = Number(id);
            const moviesResponse = await getMovieById(movieId);
            setMovie(moviesResponse);
        } catch (error) {
            console.error('Error fetching movie:', error);
        }
    }, [id]);

    useEffect(() => {
        fetchMovieById();
    }, [fetchMovieById]);

    useEffect(() => {
        if (movie) {
            const sectionMovie = document.querySelector<HTMLElement>(".section__movie");
            const sectionPlayer = document.querySelector<HTMLElement>(".section__player");
            const sectionComments = document.querySelector<HTMLElement>(".section__comments");

            if (sectionMovie) {
                AnimateElement(sectionMovie, "animate__fadeInLeft", 0);
            }
            if (sectionPlayer) {
                AnimateElement(sectionPlayer, "animate__fadeInLeft", 150);
            }
            if (sectionComments) {
                AnimateElement(sectionComments, "animate__fadeInLeft", 300);
            }
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
        </>
    );
};

export default GetMovie;

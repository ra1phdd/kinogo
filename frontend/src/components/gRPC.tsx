import { movies_v1 } from '@protos/movies_v1/movies_v1';
import { comments_v1 } from '@protos/comments_v1/comments_v1'
import {google} from "@/google/protobuf/timestamp.ts";

const clientMoviesV1 = new movies_v1.MoviesV1Client('http://localhost:10000');
const clientCommentsV1 = new comments_v1.CommentsV1Client('http://localhost:10000');

export interface Movies {
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

export interface Movie {
    id: number;
    title: string;
    description: string;
    country: string;
    releaseDate: number;
    timeMovie: number;
    scoreKP: number;
    scoreIMDB: number;
    poster: string;
    typeMovie: number;
    genres: string;
}

export interface BestMovie {
    id: number;
    title: string;
    releaseDate: number;
    poster: string;
}

export interface Timestamp {
    seconds: number;
    nanos: number;
}

export interface Comments {
    id: number;
    parentId: number;
    user: {
        username: string;
        photoUrl: string;
        firstName: string;
        lastName: string;
    };
    text: string;
    createdAt: Timestamp;
    updatedAt: Timestamp;
    children?: Comments[];
}

export function getMovies(limit: number, page: number): Promise<Movies[]> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesRequest();
        request.limit = limit;
        request.page = page;

        clientMoviesV1.GetMovies(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.movies) {
                resolve(response.movies);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function getMoviesByTypeMovie(limit: number, page: number, typeMovie: number): Promise<Movies[]> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesByFilterRequest();
        request.filters = new movies_v1.GetMoviesByFilterItem();
        request.filters.typeMovie = typeMovie;
        request.limit = limit;
        request.page = page;

        clientMoviesV1.GetMoviesByFilter(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.movies) {
                resolve(response.movies);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function getMovieById(id: number): Promise<Movie> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesByIdRequest();
        request.id = id;

        clientMoviesV1.GetMovieById(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response) {
                resolve(response);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function getBestMovie(): Promise<BestMovie[]> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesByFilterRequest();
        request.filters = new movies_v1.GetMoviesByFilterItem();
        request.filters.bestMovie = true;
        request.limit = 1;
        request.page = 1;

        clientMoviesV1.GetMoviesByFilter(request, {}, (err, response) => {
            if (err) {
                reject(err);
            }  else if (response && response.movies) {
                const bestMovies: BestMovie[] = response.movies.map(movie => ({
                    id: movie.id,
                    title: movie.title,
                    releaseDate: movie.releaseDate,
                    poster: movie.poster
                }));
                resolve(bestMovies);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function getComments(movieId: number, limit: number, page: number): Promise<Comments[]> {
    return new Promise((resolve, reject) => {
        const request = new comments_v1.GetCommentsByIdRequest();
        request.movieId = movieId;
        request.limit = limit;
        request.page = page;

        clientCommentsV1.GetCommentsById(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.comments) {
                resolve(response.comments);
            } else {
                reject(new Error('No comments found'));
            }
        });
    });
}

export function addComment(parentId: number | null, movieId: number, userId: number, text: string): Promise<number> {
    return new Promise((resolve, reject) => {
        const request = new comments_v1.AddCommentRequest();
        const now = new Date();
        const timestamp = new google.protobuf.Timestamp({ seconds: Math.floor(now.getTime() / 1000), nanos: now.getMilliseconds() * 1e6 });

        request.parentId = parentId ?? 0; // Set to 0 if parentId is null
        request.movieId = movieId;
        request.userId = userId;
        request.text = text;
        request.createdAt = timestamp;

        clientCommentsV1.AddComment(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.err == "") {
                resolve(response.id);
            } else {
                reject(new Error('Failed to add comment'));
            }
        });
    });
}

export function updateComment(id: number, text: string): Promise<string> {
    return new Promise((resolve, reject) => {
        const request = new comments_v1.UpdateCommentRequest();
        const now = new Date();
        const timestamp = new google.protobuf.Timestamp({ seconds: Math.floor(now.getTime() / 1000), nanos: now.getMilliseconds() * 1e6 });

        request.id = id;
        request.text = text;
        request.updatedAt = timestamp;

        clientCommentsV1.UpdateComment(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.err == "") {
                resolve("");
            } else {
                reject(new Error('Failed to add comment'));
            }
        });
    });
}

export function deleteComment(id: number): Promise<string> {
    return new Promise((resolve, reject) => {
        const request = new comments_v1.DelCommentRequest();
        request.id = id;

        clientCommentsV1.DelComment(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else if (response && response.err == "") {
                resolve("");
            } else {
                reject(new Error('Failed to add comment'));
            }
        });
    });
}
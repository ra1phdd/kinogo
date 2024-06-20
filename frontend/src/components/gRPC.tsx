import { movies_v1 } from '@protos/movies_v1/movies_v1';

const client = new movies_v1.MoviesV1Client('http://localhost:10000');

export function getMovies(limit: number, page: number): Promise<any> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesRequest();
        request.limit = limit;
        request.page = page;

        client.GetMovies(request, {}, (err, response) => {
            if (err) {
                console.error('Error fetching movies:', err);
                reject(err);
            } else {
                resolve(response.movies);
            }
        });
    });
}

export function getMovieById(id: number): Promise<any> {
    return new Promise((resolve, reject) => {
        const request = new movies_v1.GetMoviesByIdRequest();
        request.id = id;

        client.GetMovieById(request, {}, (err, response) => {
            if (err) {
                console.error('Error fetching movie:', err);
                reject(err);
            } else {
                resolve(response);
            }
        });
    });
}

/*export const getMovieById = (movieId: number): Promise<GetMoviesByIdResponse> => {
    const request = new GetMoviesByIdRequest();
    request.setId(movieId);

    return new Promise((resolve, reject) => {
        client.getMovieById(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else {
                resolve(response);
            }
        });
    });
};*/
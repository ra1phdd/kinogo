import * as jspb from 'google-protobuf'



export class GetMoviesRequest extends jspb.Message {
  getLimit(): number;
  setLimit(value: number): GetMoviesRequest;

  getPage(): number;
  setPage(value: number): GetMoviesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesRequest): GetMoviesRequest.AsObject;
  static serializeBinaryToWriter(message: GetMoviesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesRequest;
  static deserializeBinaryFromReader(message: GetMoviesRequest, reader: jspb.BinaryReader): GetMoviesRequest;
}

export namespace GetMoviesRequest {
  export type AsObject = {
    limit: number,
    page: number,
  }
}

export class GetMoviesByIdRequest extends jspb.Message {
  getId(): number;
  setId(value: number): GetMoviesByIdRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesByIdRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesByIdRequest): GetMoviesByIdRequest.AsObject;
  static serializeBinaryToWriter(message: GetMoviesByIdRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesByIdRequest;
  static deserializeBinaryFromReader(message: GetMoviesByIdRequest, reader: jspb.BinaryReader): GetMoviesByIdRequest;
}

export namespace GetMoviesByIdRequest {
  export type AsObject = {
    id: number,
  }
}

export class GetMoviesByFilterRequest extends jspb.Message {
  getFilters(): GetMoviesByFilterItem | undefined;
  setFilters(value?: GetMoviesByFilterItem): GetMoviesByFilterRequest;
  hasFilters(): boolean;
  clearFilters(): GetMoviesByFilterRequest;

  getLimit(): number;
  setLimit(value: number): GetMoviesByFilterRequest;

  getPage(): number;
  setPage(value: number): GetMoviesByFilterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesByFilterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesByFilterRequest): GetMoviesByFilterRequest.AsObject;
  static serializeBinaryToWriter(message: GetMoviesByFilterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesByFilterRequest;
  static deserializeBinaryFromReader(message: GetMoviesByFilterRequest, reader: jspb.BinaryReader): GetMoviesByFilterRequest;
}

export namespace GetMoviesByFilterRequest {
  export type AsObject = {
    filters?: GetMoviesByFilterItem.AsObject,
    limit: number,
    page: number,
  }
}

export class GetMoviesByFilterItem extends jspb.Message {
  getTypemovie(): number;
  setTypemovie(value: number): GetMoviesByFilterItem;

  getSearch(): string;
  setSearch(value: string): GetMoviesByFilterItem;

  getGenresList(): Array<Genres>;
  setGenresList(value: Array<Genres>): GetMoviesByFilterItem;
  clearGenresList(): GetMoviesByFilterItem;
  addGenres(value?: Genres, index?: number): Genres;

  getYearmin(): number;
  setYearmin(value: number): GetMoviesByFilterItem;

  getYearmax(): number;
  setYearmax(value: number): GetMoviesByFilterItem;

  getBestmovie(): boolean;
  setBestmovie(value: boolean): GetMoviesByFilterItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesByFilterItem.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesByFilterItem): GetMoviesByFilterItem.AsObject;
  static serializeBinaryToWriter(message: GetMoviesByFilterItem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesByFilterItem;
  static deserializeBinaryFromReader(message: GetMoviesByFilterItem, reader: jspb.BinaryReader): GetMoviesByFilterItem;
}

export namespace GetMoviesByFilterItem {
  export type AsObject = {
    typemovie: number,
    search: string,
    genresList: Array<Genres.AsObject>,
    yearmin: number,
    yearmax: number,
    bestmovie: boolean,
  }
}

export class AddMoviesRequest extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): AddMoviesRequest;

  getDescription(): string;
  setDescription(value: string): AddMoviesRequest;

  getCountriesList(): Array<Countries>;
  setCountriesList(value: Array<Countries>): AddMoviesRequest;
  clearCountriesList(): AddMoviesRequest;
  addCountries(value?: Countries, index?: number): Countries;

  getReleasedate(): number;
  setReleasedate(value: number): AddMoviesRequest;

  getTimemovie(): number;
  setTimemovie(value: number): AddMoviesRequest;

  getScorekp(): number;
  setScorekp(value: number): AddMoviesRequest;

  getScoreimdb(): number;
  setScoreimdb(value: number): AddMoviesRequest;

  getPoster(): string;
  setPoster(value: string): AddMoviesRequest;

  getTypemovie(): number;
  setTypemovie(value: number): AddMoviesRequest;

  getGenresList(): Array<Genres>;
  setGenresList(value: Array<Genres>): AddMoviesRequest;
  clearGenresList(): AddMoviesRequest;
  addGenres(value?: Genres, index?: number): Genres;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMoviesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMoviesRequest): AddMoviesRequest.AsObject;
  static serializeBinaryToWriter(message: AddMoviesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMoviesRequest;
  static deserializeBinaryFromReader(message: AddMoviesRequest, reader: jspb.BinaryReader): AddMoviesRequest;
}

export namespace AddMoviesRequest {
  export type AsObject = {
    title: string,
    description: string,
    countriesList: Array<Countries.AsObject>,
    releasedate: number,
    timemovie: number,
    scorekp: number,
    scoreimdb: number,
    poster: string,
    typemovie: number,
    genresList: Array<Genres.AsObject>,
  }
}

export class Countries extends jspb.Message {
  getName(): string;
  setName(value: string): Countries;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Countries.AsObject;
  static toObject(includeInstance: boolean, msg: Countries): Countries.AsObject;
  static serializeBinaryToWriter(message: Countries, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Countries;
  static deserializeBinaryFromReader(message: Countries, reader: jspb.BinaryReader): Countries;
}

export namespace Countries {
  export type AsObject = {
    name: string,
  }
}

export class Genres extends jspb.Message {
  getName(): string;
  setName(value: string): Genres;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Genres.AsObject;
  static toObject(includeInstance: boolean, msg: Genres): Genres.AsObject;
  static serializeBinaryToWriter(message: Genres, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Genres;
  static deserializeBinaryFromReader(message: Genres, reader: jspb.BinaryReader): Genres;
}

export namespace Genres {
  export type AsObject = {
    name: string,
  }
}

export class DeleteMoviesRequest extends jspb.Message {
  getId(): number;
  setId(value: number): DeleteMoviesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteMoviesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteMoviesRequest): DeleteMoviesRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteMoviesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteMoviesRequest;
  static deserializeBinaryFromReader(message: DeleteMoviesRequest, reader: jspb.BinaryReader): DeleteMoviesRequest;
}

export namespace DeleteMoviesRequest {
  export type AsObject = {
    id: number,
  }
}

export class GetMoviesResponse extends jspb.Message {
  getMoviesList(): Array<GetMovieItem>;
  setMoviesList(value: Array<GetMovieItem>): GetMoviesResponse;
  clearMoviesList(): GetMoviesResponse;
  addMovies(value?: GetMovieItem, index?: number): GetMovieItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesResponse): GetMoviesResponse.AsObject;
  static serializeBinaryToWriter(message: GetMoviesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesResponse;
  static deserializeBinaryFromReader(message: GetMoviesResponse, reader: jspb.BinaryReader): GetMoviesResponse;
}

export namespace GetMoviesResponse {
  export type AsObject = {
    moviesList: Array<GetMovieItem.AsObject>,
  }
}

export class GetMovieItem extends jspb.Message {
  getId(): number;
  setId(value: number): GetMovieItem;

  getTitle(): string;
  setTitle(value: string): GetMovieItem;

  getDescription(): string;
  setDescription(value: string): GetMovieItem;

  getReleasedate(): number;
  setReleasedate(value: number): GetMovieItem;

  getScorekp(): number;
  setScorekp(value: number): GetMovieItem;

  getScoreimdb(): number;
  setScoreimdb(value: number): GetMovieItem;

  getPoster(): string;
  setPoster(value: string): GetMovieItem;

  getTypemovie(): number;
  setTypemovie(value: number): GetMovieItem;

  getGenres(): string;
  setGenres(value: string): GetMovieItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMovieItem.AsObject;
  static toObject(includeInstance: boolean, msg: GetMovieItem): GetMovieItem.AsObject;
  static serializeBinaryToWriter(message: GetMovieItem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMovieItem;
  static deserializeBinaryFromReader(message: GetMovieItem, reader: jspb.BinaryReader): GetMovieItem;
}

export namespace GetMovieItem {
  export type AsObject = {
    id: number,
    title: string,
    description: string,
    releasedate: number,
    scorekp: number,
    scoreimdb: number,
    poster: string,
    typemovie: number,
    genres: string,
  }
}

export class GetMoviesByIdResponse extends jspb.Message {
  getId(): number;
  setId(value: number): GetMoviesByIdResponse;

  getTitle(): string;
  setTitle(value: string): GetMoviesByIdResponse;

  getDescription(): string;
  setDescription(value: string): GetMoviesByIdResponse;

  getCountry(): string;
  setCountry(value: string): GetMoviesByIdResponse;

  getReleasedate(): number;
  setReleasedate(value: number): GetMoviesByIdResponse;

  getTimemovie(): number;
  setTimemovie(value: number): GetMoviesByIdResponse;

  getScorekp(): number;
  setScorekp(value: number): GetMoviesByIdResponse;

  getScoreimdb(): number;
  setScoreimdb(value: number): GetMoviesByIdResponse;

  getPoster(): string;
  setPoster(value: string): GetMoviesByIdResponse;

  getTypemovie(): number;
  setTypemovie(value: number): GetMoviesByIdResponse;

  getGenres(): string;
  setGenres(value: string): GetMoviesByIdResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMoviesByIdResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMoviesByIdResponse): GetMoviesByIdResponse.AsObject;
  static serializeBinaryToWriter(message: GetMoviesByIdResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMoviesByIdResponse;
  static deserializeBinaryFromReader(message: GetMoviesByIdResponse, reader: jspb.BinaryReader): GetMoviesByIdResponse;
}

export namespace GetMoviesByIdResponse {
  export type AsObject = {
    id: number,
    title: string,
    description: string,
    country: string,
    releasedate: number,
    timemovie: number,
    scorekp: number,
    scoreimdb: number,
    poster: string,
    typemovie: number,
    genres: string,
  }
}

export class AddMoviesResponse extends jspb.Message {
  getId(): number;
  setId(value: number): AddMoviesResponse;

  getErr(): string;
  setErr(value: string): AddMoviesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMoviesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMoviesResponse): AddMoviesResponse.AsObject;
  static serializeBinaryToWriter(message: AddMoviesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMoviesResponse;
  static deserializeBinaryFromReader(message: AddMoviesResponse, reader: jspb.BinaryReader): AddMoviesResponse;
}

export namespace AddMoviesResponse {
  export type AsObject = {
    id: number,
    err: string,
  }
}

export class DeleteMoviesResponse extends jspb.Message {
  getErr(): string;
  setErr(value: string): DeleteMoviesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteMoviesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteMoviesResponse): DeleteMoviesResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteMoviesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteMoviesResponse;
  static deserializeBinaryFromReader(message: DeleteMoviesResponse, reader: jspb.BinaryReader): DeleteMoviesResponse;
}

export namespace DeleteMoviesResponse {
  export type AsObject = {
    err: string,
  }
}


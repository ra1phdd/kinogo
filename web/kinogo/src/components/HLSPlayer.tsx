import '@vidstack/react/player/styles/default/theme.css';
import '@vidstack/react/player/styles/default/layouts/video.css';
import { MediaPlayer, MediaProvider } from '@vidstack/react';
import { defaultLayoutIcons, DefaultVideoLayout } from '@vidstack/react/player/layouts/default';

export const VideoPlayer = () => {
    return (
        <MediaPlayer title="Sprite Fight" src="http://localhost:4000/api/v1/stream/7/1080p/stream.m3u8">
            <MediaProvider />
            <DefaultVideoLayout icons={defaultLayoutIcons} />
        </MediaPlayer>
    );
}

export default VideoPlayer;
import React, { useState, useRef, useCallback, useEffect } from 'react';
import {MediaPlayer, MediaPlayerInstance, MediaProvider} from '@vidstack/react';
import { defaultLayoutIcons, DefaultVideoLayout } from '@vidstack/react/player/layouts/default';
import '@vidstack/react/player/styles/default/theme.css';
import '@vidstack/react/player/styles/default/layouts/video.css';
import {metricStreamingPerformance} from "@components/gRPC.tsx";

const VideoPlayer: React.FC<{ title: string; id: number }> = ({ title, id }) => {
    const qualities = [
        { src: `http://localhost:4000/stream/${id}/1080p/playlist.m3u8`, quality: '1080p' },
        { src: `http://localhost:4000/stream/${id}/720p/playlist.m3u8`, quality: '720p' },
        { src: `http://localhost:4000/stream/${id}/480p/playlist.m3u8`, quality: '480p' },
        { src: `http://localhost:4000/stream/${id}/360p/playlist.m3u8`, quality: '360p' },
    ];

    const playerRef = useRef<MediaPlayerInstance>(null);
    const [currentQuality, setCurrentQuality] = useState(qualities[0]);
    const [lastTime, setLastTime] = useState(0);
    const [metrics, setMetrics] = useState<{
        bufferingCount: number;
        bufferingTime: number;
        startTime: number;
        playbackSuccess: boolean;
        playbackError: MediaError | null;
        viewTime: number,
        lastPlayTime: number,
        currentTime: number;
        duration: number;
        volume: number;
        isMuted: boolean;
    }>({
        bufferingCount: 0,
        bufferingTime: 0,
        startTime: 0,
        playbackSuccess: true,
        playbackError: null,
        viewTime: 0,
        lastPlayTime: 0,
        currentTime: 0,
        duration: 0,
        volume: 1,
        isMuted: false,
    });

    useEffect(() => {
        const player = playerRef.current;

        if (player) {
            // Event handlers
            const handlePlay = () => {
                setMetrics((prev) => ({
                    ...prev,
                    startTime: Date.now()
                }));
                updateMetrics();
            };

            const updateMetrics = () => {
                setMetrics((prev) => ({
                    ...prev,
                    currentTime: player.currentTime,
                    duration: player.duration,
                    volume: player.volume,
                    isMuted: player.muted,
                }));
            };

            const handlePlaying = () => {
                const currentTime = Date.now();
                const startTime = metrics.startTime;

                if (startTime) {
                    const bufferingTime = currentTime - startTime;
                    setMetrics((prev) => ({
                        ...prev,
                        bufferingTime: prev.bufferingTime + bufferingTime,
                        startTime: 0,
                    }));
                }

                setMetrics((prev) => ({
                    ...prev,
                    lastPlayTime: player.currentTime,
                }));

                updateMetrics();
            };

            const handlePause = () => {
                const viewTime = player.currentTime - metrics.lastPlayTime;
                setMetrics((prev) => ({
                    ...prev,
                    viewTime: viewTime,
                }));

                updateMetrics();

                metricStreamingPerformance(id, metrics.bufferingCount, metrics.bufferingTime, String(metrics.playbackError), metrics.viewTime, metrics.duration);
            };

            const handleBufferingStart = () => {
                setMetrics((prev) => ({ ...prev, startTime: Date.now() }));
            };

            const handleBufferingEnd = () => {
                const currentTime = Date.now();
                const startTime = metrics.startTime;

                if (startTime) {
                    const bufferingTime = currentTime - startTime;
                    setMetrics((prev) => ({
                        ...prev,
                        bufferingTime: prev.bufferingTime + bufferingTime,
                        startTime: 0,
                    }));
                }

                updateMetrics();
            };

            const handleBuffering = () => {
                setMetrics((prev) => ({ ...prev, bufferingCount: prev.bufferingCount + 1 }));
            };

            const handleError = (event: Event) => {
                console.log('Playback error', event);
                setMetrics((prev) => ({
                    ...prev,
                    playbackSuccess: false,
                    playbackError: (event.target as HTMLMediaElement).error
                }));

                metricStreamingPerformance(id, metrics.bufferingCount, metrics.bufferingTime, String(metrics.playbackError), metrics.viewTime, metrics.duration);
            };

            // Attach event listeners
            player.addEventListener('play', handlePlay);
            player.addEventListener('playing', handlePlaying);
            player.addEventListener('waiting', handleBuffering);
            player.addEventListener('error', handleError);
            player.addEventListener('waiting', handleBufferingStart);
            player.addEventListener('playing', handleBufferingEnd);
            player.addEventListener('pause', handlePause);

            return () => {
                player.removeEventListener('play', handlePlay);
                player.removeEventListener('playing', handlePlaying);
                player.removeEventListener('waiting', handleBuffering);
                player.removeEventListener('error', handleError);
                player.removeEventListener('waiting', handleBufferingStart);
                player.removeEventListener('playing', handleBufferingEnd);
                player.removeEventListener('pause', handlePause);
            };
        }
    }, [metrics]);

    useEffect(() => {
        const savedTime = localStorage.getItem(`videoTime_${id}`);
        if (savedTime) {
            setLastTime(parseFloat(savedTime));
        }
    }, [id]);

    useEffect(() => {
        const handleBeforeUnload = () => {
            if (playerRef.current) {
                localStorage.setItem(`videoTime_${id}`, playerRef.current.currentTime.toString());
            }
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, [id]);

    const russianTranslations = {
        "Accessibility": "Доступность",
        "Announcements": "Оповещения",
        "Audio": "Аудио",
        "Auto": "Авто",
        "Boost": "Усиление",
        "Buffered": "Буферизовано",
        "Current Time": "Текущее время",
        "Duration": "Продолжительность",
        "Enter Fullscreen": "На весь экран",
        "Enter PiP": "Картинка-в-картинке",
        "Exit Fullscreen": "Выйти из полноэкранного режима",
        "Keyboard Animations": "Анимация клавиатуры",
        "Loop": "Повтор",
        "Mute": "Выключить звук",
        "Normal": "Обычный",
        "Pause": "Пауза",
        "Play": "Воспроизвести",
        "Playback": "Воспроизведение",
        "Played": "Воспроизведено",
        "Quality": "Качество",
        "Remaining Time": "Оставшееся время",
        "Seek": "Перемотка",
        "Settings": "Настройки",
        "Speed": "Скорость",
        "Unmute": "Включить звук"
    };

    const handleLoadedMetadata = useCallback(() => {
        if (playerRef.current && lastTime > 0) {
            playerRef.current.currentTime = lastTime;
            setLastTime(0);
        }
    }, [lastTime]);

    const handleQualityChange = (quality: string) => {
        const newQuality = qualities.find(q => q.quality === quality);
        if (newQuality) {
            setLastTime(playerRef.current?.currentTime || 0);
            setCurrentQuality(newQuality);
        }
    };

    return (
        <>
            <MediaPlayer
                ref={playerRef}
                title={title}
                src={currentQuality.src}
                onLoadedMetadata={handleLoadedMetadata}
                streamType="on-demand">
                <MediaProvider/>
                <DefaultVideoLayout
                    icons={defaultLayoutIcons}
                    translations={russianTranslations}
                />
                <div className="quality-selector-container">
                    <div className="quality-selector">
                        {qualities.map((q) => (
                            <button
                                key={q.quality}
                                className={`quality-button ${currentQuality.quality === q.quality ? 'active' : ''}`}
                                onClick={() => handleQualityChange(q.quality)}
                            >
                                {q.quality}
                            </button>
                        ))}
                    </div>
                </div>
            </MediaPlayer>
        </>
    );
};

export default VideoPlayer;
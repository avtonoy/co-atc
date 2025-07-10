package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the main application configuration structure
// containing all configuration sections
type Config struct {
	Server         ServerConfig         `toml:"server"`          // HTTP server settings
	ADSB           ADSBConfig           `toml:"adsb"`            // Aircraft tracking data source settings
	Frequencies    FrequenciesConfig    `toml:"frequencies"`     // Radio frequency monitoring settings
	Logging        LoggingConfig        `toml:"logging"`         // Application logging settings
	Storage        StorageConfig        `toml:"storage"`         // Data persistence settings
	Station        StationConfig        `toml:"station"`         // Physical location settings
	Transcription  TranscriptionConfig  `toml:"transcription"`   // Audio transcription settings
	PostProcessing PostProcessingConfig `toml:"post_processing"` // Post-processing settings for transcriptions
	FlightPhases   FlightPhasesConfig   `toml:"flight_phases"`   // Flight phase detection settings
	Weather        WeatherConfig        `toml:"wx"`              // Weather data fetching and caching settings
	ATCChat        ATCChatConfig        `toml:"atc_chat"`        // ATC Chat voice assistant settings
	Templating     TemplatingConfig     `toml:"templating"`      // Shared templating system settings
}

// ServerConfig contains HTTP server configuration settings
type ServerConfig struct {
	Port               int      `toml:"port"`                  // Primary HTTP port for the server
	Host               string   `toml:"host"`                  // Host address to bind to (e.g., 127.0.0.1 for localhost only, 0.0.0.0 for all interfaces)
	CORSAllowedOrigins []string `toml:"cors_allowed_origins"`  // List of origins allowed for CORS requests (use ["*"] for all origins)
	ReadTimeoutSecs    int      `toml:"read_timeout_seconds"`  // Maximum duration for reading the entire request (0 = no timeout)
	WriteTimeoutSecs   int      `toml:"write_timeout_seconds"` // Maximum duration for writing the response (0 = no timeout, recommended for streaming)
	IdleTimeoutSecs    int      `toml:"idle_timeout_seconds"`  // Maximum duration to wait for the next request when keep-alives are enabled
	AdditionalPorts    []int    `toml:"additional_ports"`      // Additional HTTP ports to listen on (useful for multiple interfaces)
	StaticFilesDir     string   `toml:"static_files_dir"`      // Directory to serve static files from (e.g., "www")
}

// ADSBConfig contains ADS-B aircraft tracking data source configuration
type ADSBConfig struct {
	// Source selection
	SourceType string `toml:"source_type"` // Data source type: "local" (e.g., dump1090) or "external" (e.g., ADS-B Exchange API)

	// Legacy field - deprecated, use LocalSourceURL instead
	SourceURL string `toml:"source_url"` // DEPRECATED: Legacy URL field for backward compatibility

	// Local source settings (used when source_type = "local")
	LocalSourceURL string `toml:"local_source_url"` // URL for local ADS-B source (e.g., http://192.168.1.10/tar1090/data/aircraft.json)

	// External API source settings (used when source_type = "external")
	ExternalSourceURL string `toml:"external_source_url"` // URL template for external API with format placeholders for lat, lon, and distance
	APIHost           string `toml:"api_host"`            // API host header value (e.g., for RapidAPI)
	APIKey            string `toml:"api_key"`             // API key for authentication with external service
	SearchRadiusNM    int    `toml:"search_radius_nm"`    // Search radius in nautical miles for external API queries

	// Common settings for both source types
	FetchIntervalSecs        int    `toml:"fetch_interval_seconds"`      // How often to fetch new aircraft data (in seconds)
	SignalLostTimeoutSecs    int    `toml:"signal_lost_timeout_seconds"` // Time after which aircraft is marked as signal_lost (in seconds, default: 60)
	AirlineDBPath            string `toml:"airline_db_path"`             // Path to airline database JSON file for aircraft operator lookups
	WebSocketAircraftUpdates bool   `toml:"websocket_aircraft_updates"`  // Enable WebSocket aircraft streaming (hybrid mode)
}

// LoggingConfig contains application logging configuration
type LoggingConfig struct {
	Level  string `toml:"level"`  // Log level: "debug", "info", "warn", or "error"
	Format string `toml:"format"` // Log format: "json" (structured) or "console" (human-readable)
}

// StorageConfig contains data persistence configuration
type StorageConfig struct {
	Type              string `toml:"type"`                 // Storage backend type (currently only "sqlite" is supported)
	SQLiteBasePath    string `toml:"sqlite_base_path"`     // Base path for SQLite database files (actual filename will be generated as co-atc-YYYY-MM-DD.db)
	MaxPositionsInAPI int    `toml:"max_positions_in_api"` // Maximum number of positions to return in the /aircraft API response
}

// StationConfig contains physical location configuration for the monitoring station
type StationConfig struct {
	Latitude                float64 `toml:"latitude"`                   // Latitude of the station in decimal degrees (-90 to 90)
	Longitude               float64 `toml:"longitude"`                  // Longitude of the station in decimal degrees (-180 to 180)
	ElevationFeet           int     `toml:"elevation_feet"`             // Elevation of the station above sea level in feet
	AirportCode             string  `toml:"airport_code"`               // ICAO code of the airport (e.g., "CYYZ")
	RunwaysDBPath           string  `toml:"runways_db_path"`            // Path to runway database JSON file
	RunwayExtensionLengthNM float64 `toml:"runway_extension_length_nm"` // Length of runway extensions in nautical miles
	AirportRangeNM          float64 `toml:"airport_range_nm"`           // Range in nautical miles to consider aircraft as being at this airport (default: 5.0)
}

// TranscriptionConfig contains settings for audio transcription services
type TranscriptionConfig struct {
	// OpenAI API settings
	OpenAIAPIKey string `toml:"openai_api_key"` // OpenAI API key for transcription service
	Model        string `toml:"model"`          // OpenAI model to use (e.g., "gpt-4o-transcribe")
	Language     string `toml:"language"`       // Primary language for transcription (e.g., "en" for English)
	PromptPath   string `toml:"prompt_path"`    // Path to the system prompt file for transcription

	// Audio processing settings
	NoiseReduction string `toml:"noise_reduction"` // Noise reduction mode: "near_field", "far_field", or "none"
	ChunkMs        int    `toml:"chunk_ms"`        // Size of audio chunks for processing in milliseconds
	BufferSizeKB   int    `toml:"buffer_size_kb"`  // Audio buffer size in kilobytes

	// FFmpeg conversion settings
	FFmpegPath       string `toml:"ffmpeg_path"`        // Path to FFmpeg executable
	FFmpegSampleRate int    `toml:"ffmpeg_sample_rate"` // Audio sample rate in Hz (typically 24000 for OpenAI)
	FFmpegChannels   int    `toml:"ffmpeg_channels"`    // Number of audio channels (1 for mono, 2 for stereo)
	FFmpegFormat     string `toml:"ffmpeg_format"`      // Audio format (e.g., "s16le" for signed 16-bit little-endian PCM)

	// Connection management
	ReconnectIntervalSec int `toml:"reconnect_interval_sec"` // Seconds to wait before reconnecting after failure
	MaxRetries           int `toml:"max_retries"`            // Maximum number of connection retry attempts

	// Voice activity detection (VAD) settings
	TurnDetectionType string  `toml:"turn_detection_type"` // Method for detecting speech turns (e.g., "server_vad")
	PrefixPaddingMs   int     `toml:"prefix_padding_ms"`   // Milliseconds of audio to include before detected speech
	SilenceDurationMs int     `toml:"silence_duration_ms"` // Milliseconds of silence to consider end of speech
	VADThreshold      float64 `toml:"vad_threshold"`       // Threshold for voice activity detection (0.0-1.0)

	// API retry settings
	RetryMaxAttempts      int `toml:"retry_max_attempts"`       // Maximum number of API call retry attempts
	RetryInitialBackoffMs int `toml:"retry_initial_backoff_ms"` // Initial backoff time in milliseconds
	RetryMaxBackoffMs     int `toml:"retry_max_backoff_ms"`     // Maximum backoff time in milliseconds

	// HTTP timeout settings
	TimeoutSeconds int `toml:"timeout_seconds"` // HTTP timeout for OpenAI API requests in seconds
}

// PostProcessingConfig contains settings for post-processing of transcriptions
type PostProcessingConfig struct {
	Enabled               bool   `toml:"enabled"`                // Enable or disable post-processing
	Model                 string `toml:"model"`                  // OpenAI model to use for post-processing
	IntervalSeconds       int    `toml:"interval_seconds"`       // How often to run the post-processing (in seconds)
	BatchSize             int    `toml:"batch_size"`             // Maximum number of transcriptions to process in each batch
	ContextTranscriptions int    `toml:"context_transcriptions"` // Number of previous processed transcriptions to include for context
	SystemPromptPath      string `toml:"system_prompt_path"`     // Path to the system prompt file
	TimeoutSeconds        int    `toml:"timeout_seconds"`        // HTTP timeout for OpenAI API requests in seconds
}

// FrequenciesConfig contains settings for radio frequency monitoring
type FrequenciesConfig struct {
	Sources               []FrequencyConfig `toml:"sources"`                 // List of radio frequencies to monitor
	BufferSizeKB          int               `toml:"buffer_size_kb"`          // Audio buffer size in kilobytes
	StreamTimeoutSecs     int               `toml:"stream_timeout_secs"`     // Timeout for audio streams (0 = no timeout)
	ReconnectIntervalSecs int               `toml:"reconnect_interval_secs"` // Seconds to wait before reconnecting after stream failure

	// FFmpeg timeout configuration
	FFmpegTimeoutSecs        int `toml:"ffmpeg_timeout_secs"`         // FFmpeg connection timeout in seconds (0 = no timeout, default: 30)
	FFmpegReconnectDelaySecs int `toml:"ffmpeg_reconnect_delay_secs"` // FFmpeg reconnect delay in seconds (default: 2)
}

// FrequencyConfig contains configuration for a single monitored radio frequency
type FrequencyConfig struct {
	ID              string  `toml:"id"`               // Unique identifier for this frequency
	Airport         string  `toml:"airport"`          // ICAO code of the airport (e.g., "CYYZ" for Toronto Pearson)
	Name            string  `toml:"name"`             // Human-readable name (e.g., "CYYZ Tower")
	FrequencyMHz    float64 `toml:"frequency_mhz"`    // Actual radio frequency in MHz (e.g., 118.7)
	URL             string  `toml:"url"`              // URL to the audio stream
	Order           int     `toml:"order"`            // Display order in the UI (lower numbers first)
	TranscribeAudio bool    `toml:"transcribe_audio"` // Whether to transcribe audio for this frequency
}

// FlightPhasesConfig contains settings for flight phase detection
type FlightPhasesConfig struct {
	Enabled                       bool    `toml:"enabled"`                          // Enable enhanced flight phase detection
	CruiseAltitudeFt              int     `toml:"cruise_altitude_ft"`               // Minimum altitude for cruise phase
	DepartureAltitudeFt           int     `toml:"departure_altitude_ft"`            // Minimum altitude for departure phase
	TaxiingMinSpeedKts            int     `toml:"taxiing_min_speed_kts"`            // Minimum ground speed for taxiing
	TaxiingMaxSpeedKts            int     `toml:"taxiing_max_speed_kts"`            // Maximum ground speed for taxiing
	ApproachCenterlineToleranceNM float64 `toml:"approach_centerline_tolerance_nm"` // Distance from runway centerline
	ApproachMaxDistanceNM         int     `toml:"approach_max_distance_nm"`         // Maximum distance from runway threshold
	ApproachHeadingToleranceDeg   float64 `toml:"approach_heading_tolerance_deg"`   // Heading alignment tolerance

	// Timeout configurations - grouped together for clarity
	// These control different aspects of phase transition timing:

	// 1. Long-term inactive aircraft cleanup (default: 3600 seconds = 1 hour)
	// Aircraft that haven't had ANY phase change for this duration revert to NEW
	// This is the ultimate cleanup for parked/inactive aircraft
	PhaseChangeTimeoutSeconds int `toml:"phase_change_timeout_seconds"`

	// 2. Ground phase transition prevention (default: 60 seconds, but you may want 1800 = 30 min)
	// Prevents TAX→NEW and T/D→NEW transitions when aircraft briefly stops
	// Helps maintain phase stability during normal ground operations
	PhaseTransitionTimeoutSeconds int `toml:"phase_transition_timeout_seconds"`

	// 3. Critical phase preservation (default: 60 seconds)
	// Keeps T/O and T/D phases visible for at least this duration
	// Ensures pilots/controllers can see these important events
	PhasePreservationSeconds int `toml:"phase_preservation_seconds"`

	// 4. Phase flapping prevention (NEW - default: 300 seconds = 5 minutes)
	// Prevents rapid transitions between DEP↔APP, ARR↔DEP, and T/O↔T/D
	// Much shorter than phase_change_timeout as these are active aircraft
	PhaseFlappingPreventionSeconds int `toml:"phase_flapping_prevention_seconds"`

	// 5. Recent takeoff detection window (default: 30 minutes)
	// How long after takeoff an aircraft is considered "recently departed"
	// Used for DEP phase eligibility
	RecentTakeoffTimeoutMinutes int `toml:"recent_takeoff_timeout_minutes"`

	// Other phase detection parameters
	AirportRangeNM                   float64  `toml:"airport_range_nm"`                     // Distance considered "close to airport"
	ClimbingVerticalRateFPM          int      `toml:"climbing_vertical_rate_fpm"`           // Minimum vertical rate for climbing
	TakeoffAltitudeThresholdFt       int      `toml:"takeoff_altitude_threshold_ft"`        // Altitude threshold for takeoff detection
	EmergencySquawkCodes             []string `toml:"emergency_squawk_codes"`               // List of emergency squawk codes
	ApproachVerticalRateThresholdFPM int      `toml:"approach_vertical_rate_threshold_fpm"` // Maximum vertical rate for approach

	// Ground detection thresholds (NEW - making existing constants configurable)
	FlyingMinTASKts         float64 `toml:"flying_min_tas_kts"`        // Minimum true airspeed to be considered flying
	FlyingMinAltFt          float64 `toml:"flying_min_alt_ft"`         // Minimum altitude to be considered flying
	HelicopterAltMultiplier float64 `toml:"helicopter_alt_multiplier"` // Multiplier for helicopter detection
	HighSpeedThresholdKts   float64 `toml:"high_speed_threshold_kts"`  // Speed above which aircraft is always considered flying
	HighAltitudeOverrideFt  float64 `toml:"high_altitude_override_ft"` // Altitude above which aircraft is always considered flying (handles bad speed data)

	// Enhanced sensor validation thresholds (NEW)
	ImpossibleAltDropThresholdFt    float64 `toml:"impossible_alt_drop_threshold_ft"`    // Altitude above which drops to zero are always considered sensor errors
	ImpossibleSpeedDropThresholdKts float64 `toml:"impossible_speed_drop_threshold_kts"` // Speed above which drops to zero at altitude are considered sensor errors
	ImpossibleSpeedDropMinAltFt     float64 `toml:"impossible_speed_drop_min_alt_ft"`    // Minimum altitude for impossible speed drop detection

	// Signal lost landing detection (NEW)
	SignalLostLandingEnabled  bool    `toml:"signal_lost_landing_enabled"`    // Enable automatic landing detection for signal lost aircraft
	SignalLostLandingMaxAltFt float64 `toml:"signal_lost_landing_max_alt_ft"` // Max altitude for signal lost landing detection
}

// Load loads the configuration from the specified file path
func Load(path string) (*Config, error) {
	var config Config

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	// Read the config file
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}

// LoadWithFallback loads the configuration by checking multiple locations in order of preference
func LoadWithFallback(preferredPath string) (*Config, error) {
	// List of paths to check in order of preference
	searchPaths := []string{
		preferredPath,         // User-specified path (if provided)
		"configs/config.toml", // Legacy location in configs/ folder
		"config.toml",         // Root directory
	}

	// Remove duplicates while preserving order
	uniquePaths := make([]string, 0, len(searchPaths))
	seen := make(map[string]bool)
	for _, path := range searchPaths {
		if path != "" && !seen[path] {
			uniquePaths = append(uniquePaths, path)
			seen[path] = true
		}
	}

	var lastErr error
	for _, path := range uniquePaths {
		if _, err := os.Stat(path); err == nil {
			// File exists, try to load it
			config, err := Load(path)
			if err != nil {
				lastErr = fmt.Errorf("failed to load config from %s: %w", path, err)
				continue
			}
			return config, nil
		}
		lastErr = fmt.Errorf("config file not found: %s", path)
	}

	return nil, fmt.Errorf("config file not found in any of the expected locations: %v. Last error: %w", uniquePaths, lastErr)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate frequencies config
	if err := c.ValidateFrequencies(); err != nil {
		return err
	}

	// Validate post-processing config
	if c.PostProcessing.Enabled && c.PostProcessing.ContextTranscriptions < 0 {
		return fmt.Errorf("invalid context_transcriptions value: %d (must be >= 0)", c.PostProcessing.ContextTranscriptions)
	}

	// Validate server config
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	// Validate AdditionalPorts
	portsSeen := make(map[int]bool)
	portsSeen[c.Server.Port] = true
	for _, p := range c.Server.AdditionalPorts {
		if p <= 0 || p > 65535 {
			return fmt.Errorf("invalid additional server port: %d", p)
		}
		if portsSeen[p] {
			return fmt.Errorf("duplicate port configured: %d (primary or additional)", p)
		}
		portsSeen[p] = true
	}

	// Set default static files directory if not specified
	if c.Server.StaticFilesDir == "" {
		c.Server.StaticFilesDir = "www"
	}

	// Validate static files directory exists
	if _, err := os.Stat(c.Server.StaticFilesDir); os.IsNotExist(err) {
		return fmt.Errorf("static files directory does not exist: %s", c.Server.StaticFilesDir)
	}

	// Validate ADSB config
	if c.ADSB.SourceType == "" {
		c.ADSB.SourceType = "local" // Default to local if not specified
	}

	if c.ADSB.SourceType != "local" && c.ADSB.SourceType != "external" {
		return fmt.Errorf("invalid ADSB source type: %s (must be 'local' or 'external')", c.ADSB.SourceType)
	}

	// Handle legacy configuration
	if c.ADSB.SourceType == "local" && c.ADSB.LocalSourceURL == "" && c.ADSB.SourceURL != "" {
		c.ADSB.LocalSourceURL = c.ADSB.SourceURL
	}

	// Validate source URL based on source type
	if c.ADSB.SourceType == "local" && c.ADSB.LocalSourceURL == "" {
		return fmt.Errorf("local_source_url is required when source_type is local")
	}

	if c.ADSB.SourceType == "external" {
		if c.ADSB.ExternalSourceURL == "" {
			return fmt.Errorf("external_source_url is required when source_type is external")
		}
		if c.ADSB.APIHost == "" {
			return fmt.Errorf("api_host is required when source_type is external")
		}
		if c.ADSB.APIKey == "" {
			return fmt.Errorf("api_key is required when source_type is external")
		}
		if c.ADSB.SearchRadiusNM <= 0 {
			return fmt.Errorf("search_radius_nm must be positive when source_type is external")
		}
	}

	if c.ADSB.FetchIntervalSecs <= 0 {
		return fmt.Errorf("invalid fetch interval: %d", c.ADSB.FetchIntervalSecs)
	}
	// Set default value for MaxPositionsInAPI if not specified
	if c.Storage.MaxPositionsInAPI <= 0 {
		c.Storage.MaxPositionsInAPI = 60 // Default to 60 positions if not specified
	}

	// Validate logging config
	switch c.Logging.Level {
	case "debug", "info", "warn", "error":
		// Valid log level
	default:
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	switch c.Logging.Format {
	case "json", "console":
		// Valid log format
	default:
		return fmt.Errorf("invalid log format: %s", c.Logging.Format)
	}

	// Validate storage config
	if c.Storage.Type != "sqlite" {
		return fmt.Errorf("invalid storage type: %s (only 'sqlite' is supported)", c.Storage.Type)
	}

	if c.Storage.Type == "sqlite" && c.Storage.SQLiteBasePath == "" {
		return fmt.Errorf("sqlite_base_path is required when storage type is sqlite")
	}

	// Validate Station config
	if err := c.ValidateStation(); err != nil {
		return err
	}

	// Validate Flight Phases config
	if err := c.ValidateFlightPhases(); err != nil {
		return err
	}

	// Validate Weather config
	if err := c.ValidateWeather(); err != nil {
		return err
	}

	// Validate OpenAI API keys for enabled features
	if err := c.ValidateOpenAIKeys(); err != nil {
		return err
	}

	return nil
}

// ValidateStation validates the station configuration
func (c *Config) ValidateStation() error {
	// Validate Latitude
	if c.Station.Latitude < -90 || c.Station.Latitude > 90 {
		return fmt.Errorf("invalid station latitude: %f", c.Station.Latitude)
	}

	// Validate Longitude
	if c.Station.Longitude < -180 || c.Station.Longitude > 180 {
		return fmt.Errorf("invalid station longitude: %f", c.Station.Longitude)
	}

	// Elevation can be negative, so we'll just check if it's within a reasonable range, e.g. -2000 to 30000 feet.
	if c.Station.ElevationFeet < -2000 || c.Station.ElevationFeet > 30000 {
		return fmt.Errorf("station elevation out of typical range: %d ft", c.Station.ElevationFeet)
	}

	// Airport code validation is now handled in ValidateWeather method

	return nil
}

// ValidateFrequencies validates the frequencies configuration
func (c *Config) ValidateFrequencies() error {
	// Skip validation if no frequencies are configured
	if len(c.Frequencies.Sources) == 0 {
		return nil
	}

	// Validate buffer size
	if c.Frequencies.BufferSizeKB <= 0 {
		return fmt.Errorf("invalid buffer size: %d KB", c.Frequencies.BufferSizeKB)
	}

	// Validate stream timeout
	if c.Frequencies.StreamTimeoutSecs < 0 {
		return fmt.Errorf("invalid stream timeout: %d", c.Frequencies.StreamTimeoutSecs)
	}

	// Validate reconnect interval
	if c.Frequencies.ReconnectIntervalSecs <= 0 {
		return fmt.Errorf("invalid reconnect interval: %d", c.Frequencies.ReconnectIntervalSecs)
	}

	// Validate FFmpeg timeout configuration
	if c.Frequencies.FFmpegTimeoutSecs < 0 {
		return fmt.Errorf("invalid ffmpeg_timeout_secs: %d (must be >= 0)", c.Frequencies.FFmpegTimeoutSecs)
	}
	if c.Frequencies.FFmpegReconnectDelaySecs < 0 {
		return fmt.Errorf("invalid ffmpeg_reconnect_delay_secs: %d (must be >= 0)", c.Frequencies.FFmpegReconnectDelaySecs)
	}

	// Set default values for FFmpeg timeout configuration if not specified
	// FFmpegTimeoutSecs defaults to 0 (no timeout) - no need to set explicitly
	if c.Frequencies.FFmpegReconnectDelaySecs == 0 {
		c.Frequencies.FFmpegReconnectDelaySecs = 2 // Default to 2 seconds
	}

	// Validate frequency sources
	idMap := make(map[string]bool)
	orderMap := make(map[int]string) // Track orders to check for duplicates
	for i, freq := range c.Frequencies.Sources {
		// Validate ID
		if freq.ID == "" {
			return fmt.Errorf("frequency #%d: ID is required", i+1)
		}
		if idMap[freq.ID] {
			return fmt.Errorf("frequency #%d: duplicate ID: %s", i+1, freq.ID)
		}
		idMap[freq.ID] = true

		// Validate airport
		if freq.Airport == "" {
			return fmt.Errorf("frequency #%d: airport is required", i+1)
		}

		// Validate name
		if freq.Name == "" {
			return fmt.Errorf("frequency #%d: name is required", i+1)
		}

		// Validate frequency
		if freq.FrequencyMHz <= 0 {
			return fmt.Errorf("frequency #%d: invalid frequency: %f", i+1, freq.FrequencyMHz)
		}

		// Validate URL
		if freq.URL == "" {
			return fmt.Errorf("frequency #%d: URL is required", i+1)
		}

		// Validate order
		if freq.Order <= 0 {
			return fmt.Errorf("frequency #%d: order must be a positive integer", i+1)
		}
		if existingID, exists := orderMap[freq.Order]; exists {
			return fmt.Errorf("frequency #%d: duplicate order value %d (already used by %s)", i+1, freq.Order, existingID)
		}
		orderMap[freq.Order] = freq.ID
	}

	return nil
}

// ValidateFlightPhases validates the flight phases configuration
func (c *Config) ValidateFlightPhases() error {
	if !c.FlightPhases.Enabled {
		return nil // Skip validation if flight phases are disabled
	}

	// Set default values for new fields if not specified
	if c.FlightPhases.FlyingMinTASKts == 0 {
		c.FlightPhases.FlyingMinTASKts = 50.0
	}
	if c.FlightPhases.FlyingMinAltFt == 0 {
		c.FlightPhases.FlyingMinAltFt = 700.0
	}
	if c.FlightPhases.HelicopterAltMultiplier == 0 {
		c.FlightPhases.HelicopterAltMultiplier = 2.0
	}
	if c.FlightPhases.HighSpeedThresholdKts == 0 {
		c.FlightPhases.HighSpeedThresholdKts = 200.0
	}
	if c.FlightPhases.PhasePreservationSeconds == 0 {
		c.FlightPhases.PhasePreservationSeconds = 60
	}
	if c.FlightPhases.PhaseTransitionTimeoutSeconds == 0 {
		c.FlightPhases.PhaseTransitionTimeoutSeconds = 60
	}
	if c.FlightPhases.PhaseFlappingPreventionSeconds == 0 {
		c.FlightPhases.PhaseFlappingPreventionSeconds = 300 // 5 minutes default
	}
	if c.FlightPhases.SignalLostLandingMaxAltFt == 0 {
		c.FlightPhases.SignalLostLandingMaxAltFt = 1000.0
	}

	// Validate altitude thresholds
	if c.FlightPhases.CruiseAltitudeFt <= 0 {
		return fmt.Errorf("cruise_altitude_ft must be positive: %d", c.FlightPhases.CruiseAltitudeFt)
	}
	if c.FlightPhases.DepartureAltitudeFt <= 0 {
		return fmt.Errorf("departure_altitude_ft must be positive: %d", c.FlightPhases.DepartureAltitudeFt)
	}

	// Validate speed thresholds
	if c.FlightPhases.TaxiingMinSpeedKts < 0 {
		return fmt.Errorf("taxiing_min_speed_kts must be non-negative: %d", c.FlightPhases.TaxiingMinSpeedKts)
	}
	if c.FlightPhases.TaxiingMaxSpeedKts <= c.FlightPhases.TaxiingMinSpeedKts {
		return fmt.Errorf("taxiing_max_speed_kts (%d) must be greater than taxiing_min_speed_kts (%d)",
			c.FlightPhases.TaxiingMaxSpeedKts, c.FlightPhases.TaxiingMinSpeedKts)
	}

	// Validate approach detection parameters
	if c.FlightPhases.ApproachCenterlineToleranceNM <= 0 {
		return fmt.Errorf("approach_centerline_tolerance_nm must be positive: %f", c.FlightPhases.ApproachCenterlineToleranceNM)
	}
	if c.FlightPhases.ApproachMaxDistanceNM <= 0 {
		return fmt.Errorf("approach_max_distance_nm must be positive: %d", c.FlightPhases.ApproachMaxDistanceNM)
	}
	if c.FlightPhases.ApproachHeadingToleranceDeg <= 0 || c.FlightPhases.ApproachHeadingToleranceDeg > 180 {
		return fmt.Errorf("approach_heading_tolerance_deg must be between 1 and 180: %f", c.FlightPhases.ApproachHeadingToleranceDeg)
	}

	// Validate new ground detection thresholds
	if c.FlightPhases.FlyingMinTASKts <= 0 {
		return fmt.Errorf("flying_min_tas_kts must be positive: %f", c.FlightPhases.FlyingMinTASKts)
	}
	if c.FlightPhases.FlyingMinAltFt <= 0 {
		return fmt.Errorf("flying_min_alt_ft must be positive: %f", c.FlightPhases.FlyingMinAltFt)
	}
	if c.FlightPhases.HelicopterAltMultiplier <= 0 {
		return fmt.Errorf("helicopter_alt_multiplier must be positive: %f", c.FlightPhases.HelicopterAltMultiplier)
	}
	if c.FlightPhases.HighSpeedThresholdKts <= 0 {
		return fmt.Errorf("high_speed_threshold_kts must be positive: %f", c.FlightPhases.HighSpeedThresholdKts)
	}

	// Validate phase preservation times
	if c.FlightPhases.PhasePreservationSeconds <= 0 {
		return fmt.Errorf("phase_preservation_seconds must be positive: %d", c.FlightPhases.PhasePreservationSeconds)
	}
	if c.FlightPhases.PhaseTransitionTimeoutSeconds <= 0 {
		return fmt.Errorf("phase_transition_timeout_seconds must be positive: %d", c.FlightPhases.PhaseTransitionTimeoutSeconds)
	}
	if c.FlightPhases.PhaseFlappingPreventionSeconds <= 0 {
		return fmt.Errorf("phase_flapping_prevention_seconds must be positive: %d", c.FlightPhases.PhaseFlappingPreventionSeconds)
	}

	// Validate signal lost landing detection
	if c.FlightPhases.SignalLostLandingEnabled && c.FlightPhases.SignalLostLandingMaxAltFt <= 0 {
		return fmt.Errorf("signal_lost_landing_max_alt_ft must be positive when signal_lost_landing_enabled is true: %f", c.FlightPhases.SignalLostLandingMaxAltFt)
	}

	return nil
}

// ValidateOpenAIKeys validates OpenAI API keys for enabled features
func (c *Config) ValidateOpenAIKeys() error {
	// Check transcription API key - transcription is always available if configured
	if c.Transcription.OpenAIAPIKey == "" {
		fmt.Printf("WARN: No OpenAI API key provided for transcription - transcription features will be disabled\n")
	}

	// Check ATC chat API key if ATC chat is enabled
	if c.ATCChat.Enabled && c.ATCChat.OpenAIAPIKey == "" {
		fmt.Printf("WARN: ATC Chat is enabled but no OpenAI API key provided - ATC chat features will be disabled\n")
	}

	// Check post-processing API key if post-processing is enabled
	if c.PostProcessing.Enabled {
		// Post-processing uses the same API key as transcription
		if c.Transcription.OpenAIAPIKey == "" {
			fmt.Printf("WARN: Post-processing is enabled but no OpenAI API key provided in transcription config - post-processing features will be disabled\n")
		}
	}

	return nil
}

// ValidateWeather validates the weather configuration
func (c *Config) ValidateWeather() error {
	// Validate refresh interval
	if c.Weather.RefreshIntervalMinutes <= 0 {
		return fmt.Errorf("weather refresh_interval_minutes must be greater than 0: %d", c.Weather.RefreshIntervalMinutes)
	}

	// Validate request timeout
	if c.Weather.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("weather request_timeout_seconds must be greater than 0: %d", c.Weather.RequestTimeoutSeconds)
	}

	// Validate max retries
	if c.Weather.MaxRetries < 0 {
		return fmt.Errorf("weather max_retries must be 0 or greater: %d", c.Weather.MaxRetries)
	}

	// Validate cache expiry
	if c.Weather.CacheExpiryMinutes <= 0 {
		return fmt.Errorf("weather cache_expiry_minutes must be greater than 0: %d", c.Weather.CacheExpiryMinutes)
	}

	// Validate API base URL
	if c.Weather.APIBaseURL == "" {
		return fmt.Errorf("weather api_base_url cannot be empty")
	}

	// At least one weather type must be enabled
	if !c.Weather.FetchMETAR && !c.Weather.FetchTAF && !c.Weather.FetchNOTAMs {
		return fmt.Errorf("at least one weather type must be enabled (fetch_metar, fetch_taf, or fetch_notams)")
	}

	// Validate that airport code is set if weather fetching is enabled
	if (c.Weather.FetchMETAR || c.Weather.FetchTAF || c.Weather.FetchNOTAMs) && c.Station.AirportCode == "" {
		return fmt.Errorf("station airport_code is required when weather fetching is enabled")
	}

	return nil
}

// WeatherConfig contains weather data fetching and caching configuration
type WeatherConfig struct {
	RefreshIntervalMinutes int    `toml:"refresh_interval_minutes"` // Weather data refresh interval in minutes
	APIBaseURL             string `toml:"api_base_url"`             // Base URL for weather API (e.g., https://node.windy.com/airports)
	RequestTimeoutSeconds  int    `toml:"request_timeout_seconds"`  // HTTP request timeout in seconds
	MaxRetries             int    `toml:"max_retries"`              // Maximum number of retry attempts for failed requests
	FetchMETAR             bool   `toml:"fetch_metar"`              // Whether to fetch METAR data
	FetchTAF               bool   `toml:"fetch_taf"`                // Whether to fetch TAF data
	FetchNOTAMs            bool   `toml:"fetch_notams"`             // Whether to fetch NOTAM data
	CacheExpiryMinutes     int    `toml:"cache_expiry_minutes"`     // How long to keep cached data if refresh fails
}

// ATCChatConfig contains ATC Chat voice assistant configuration
type ATCChatConfig struct {
	// Feature toggle
	Enabled bool `toml:"enabled"` // Enable or disable ATC Chat feature

	// OpenAI API settings
	OpenAIAPIKey  string `toml:"openai_api_key"` // OpenAI API key for realtime chat
	RealtimeModel string `toml:"realtime_model"` // OpenAI realtime model to use
	Voice         string `toml:"voice"`          // Voice for audio responses

	// Audio settings
	InputAudioFormat  string `toml:"input_audio_format"`  // Input audio format (e.g., "pcm16")
	OutputAudioFormat string `toml:"output_audio_format"` // Output audio format (e.g., "pcm16")
	SampleRate        int    `toml:"sample_rate"`         // Audio sample rate in Hz
	Channels          int    `toml:"channels"`            // Number of audio channels

	// Session settings
	MaxResponseTokens int     `toml:"max_response_tokens"` // Maximum tokens in response
	Temperature       float64 `toml:"temperature"`         // Response randomness (0.0-1.0)
	Speed             float64 `toml:"speed"`               // Response speed (1.0-4.0)
	TurnDetectionType string  `toml:"turn_detection_type"` // Turn detection method
	VADThreshold      float64 `toml:"vad_threshold"`       // Voice activity detection threshold
	SilenceDurationMs int     `toml:"silence_duration_ms"` // Silence duration for turn detection

	// Context settings
	MaxContextAircraft          int `toml:"max_context_aircraft"`          // Maximum aircraft to include in context
	TranscriptionHistorySeconds int `toml:"transcription_history_seconds"` // Seconds of transcription history to include

	// System prompt configuration
	SystemPromptPath        string `toml:"system_prompt_path"`    // Path to system prompt template file
	RefreshSystemPromptSecs int    `toml:"refresh_system_prompt"` // Automatic system prompt refresh interval in seconds (0 = disabled)
}

// TemplatingConfig contains shared templating system configuration
type TemplatingConfig struct {
	// Feature toggle
	Enabled bool `toml:"enabled"` // Enable or disable templating system

	// Template cache settings
	TemplateCacheSize int  `toml:"template_cache_size"` // Maximum number of templates to cache
	ReloadTemplates   bool `toml:"reload_templates"`    // Whether to reload templates from disk (development mode)

	// ATC Chat template settings
	ATCChat TemplatingATCChatConfig `toml:"atc_chat"`

	// Post-processing template settings
	PostProcessing TemplatingPostProcessingConfig `toml:"post_processing"`
}

// TemplatingATCChatConfig contains ATC chat specific templating settings
type TemplatingATCChatConfig struct {
	TemplatePath string `toml:"template_path"` // Path to ATC chat template file
	MaxAircraft  int    `toml:"max_aircraft"`  // Maximum aircraft to include in template
}

// TemplatingPostProcessingConfig contains post-processing specific templating settings
type TemplatingPostProcessingConfig struct {
	TemplatePath          string `toml:"template_path"`          // Path to post-processing template file
	ContextTranscriptions int    `toml:"context_transcriptions"` // Number of context transcriptions to include
}

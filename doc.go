// Package speechnorm rewrites digit patterns in free-form text into
// locale-appropriate spoken words for TTS input. Supported locales:
// ar, de, en, es, fr, it, pt. The single entry point is NormaliseNumbers.
//
// Currency words (unit names and joiners) follow the locale's converter.
//
// The package has zero non-stdlib dependencies.
package speechnorm

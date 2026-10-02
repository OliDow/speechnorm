// Package speechnorm rewrites digit patterns in free-form text into
// locale-appropriate spoken words for TTS input. Supported locales:
// ar, de, en, es, fr, it, pt. The single entry point is NormaliseNumbers.
//
// Currency joiners ("and", "et") and unit nouns ("dollars", "euros", etc.)
// are locale-aware and supplied by each Converter implementation.
//
// The package has zero non-stdlib dependencies.
package speechnorm

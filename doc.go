// Package darwini provides web programming conveniences.
//
// The hierarchical URL path based request multiplexers manipulate
// Request.URL.Path so that only the delegated part of the path
// remains. The multiplexers operate on the first segment in
// Request.URL.Path. This allows constructing hierarchies easily.
package darwini

# Assets Folder

This folder contains static assets for the Smart Fit Girl web application.

## Folder Structure

```
public/assets/
├── videos/          # Video files (background videos, demos, etc.)
├── images/          # Image files (logos, icons, photos, etc.)
└── README.md        # This file
```

## Usage

### Videos
Place your workout background video in `videos/` folder:
- `videos/workout-background.mp4` - Main background video for home page

### Images  
Place images in `images/` folder:
- `images/logo.png` - App logo
- `images/icons/` - Custom icons
- `images/photos/` - Photography assets

## Accessing Assets in React

Since these files are in the `public` folder, reference them with absolute paths:

```jsx
// Videos
<video src="/assets/videos/workout-background.mp4" />

// Images  
<img src="/assets/images/logo.png" alt="Smart Fit Girl" />
```

## File Naming Conventions

- Use kebab-case: `workout-background.mp4`
- Be descriptive: `hero-section-video.mp4`
- Include dimensions for images when relevant: `logo-512x512.png`

## Optimization Recommendations

### Videos
- Use MP4 format for best browser compatibility
- Keep file size under 50MB for web delivery
- Consider multiple formats (MP4, WebM) for optimization
- Use appropriate compression settings

### Images
- Use WebP format when possible for better compression
- Provide fallback formats (PNG, JPG)
- Optimize file sizes for web delivery
- Use appropriate dimensions (don't serve oversized images)

###### Link to thesis' pdf
https://gitlab.lrz.de/gbs-cm/ba/kamber/-/jobs/artifacts/main/raw/build/main.pdf?job=pdf

### How to include the generated swift classes for the iOS client
1. Navigate to the api/ directory
2. Run ./generate.sh
3. Open the iOS client app and select the project file in the XCode navigator panel
4. Go to File > Add Files to "HeatmapUIKit"... and select the generated "/api" folder to add it in the project's root directory. Make sure to select "Create groups" instead of "Create folder references".
5. Alternatively to Steps 3 and 4: Drag and drop the files to the XCode editor

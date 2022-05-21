###### Link to thesis' pdf
https://gitlab.lrz.de/gbs-cm/ba/kamber/-/jobs/artifacts/main/raw/build/main.pdf?job=pdf

### How to include the generated swift classes for the iOS client
1. Run ./generate.sh
2. Open the iOS client app and select the project file in the XCode navigator panel
3. Go to File > Add Files to "HeatmapUIKit"... and select the generated "/api" folder to add it in the project's root directory. Make sure to select "Create groups" instead of "Create folder references".
4. Alternatively to Steps 2 and 3: Drag and drop the files to the XCode editor

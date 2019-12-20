VERSION=$(git describe --tag --abbrev=0)
if [[ $VERSION != $(git describe --tag) ]]; then
  PATCH=$(echo $VERSION | cut -d '.' -f 3)
  PATCH=$((PATCH+1))
  VERSION=$(echo $VERSION | sed -E "s/\.[0-9]+$/.$PATCH/")-dev
fi
echo $VERSION

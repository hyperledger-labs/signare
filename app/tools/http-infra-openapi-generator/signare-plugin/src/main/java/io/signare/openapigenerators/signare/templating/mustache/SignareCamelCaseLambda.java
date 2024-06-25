package io.signare.openapigenerators.signare.templating.mustache;

import com.github.benmanes.caffeine.cache.Cache;
import com.github.benmanes.caffeine.cache.Caffeine;
import com.github.benmanes.caffeine.cache.Ticker;
import com.samskivert.mustache.Mustache;
import com.samskivert.mustache.Template;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.openapitools.codegen.CodegenConfig;
import org.openapitools.codegen.config.GlobalSettings;

import java.io.IOException;
import java.io.Writer;
import java.util.Locale;
import java.util.concurrent.TimeUnit;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Converts text in a fragment to camelCase (with custom rules).
 *
 * Register:
 * <pre>
 * additionalProperties.put("signareCamelCase", new SignareCamelCaseLambda());
 * </pre>
 *
 * Use:
 * <pre>
 * {{#signareCamelCase}}{{name}}{{/signareCamelCase}}
 * </pre>
 */
public class SignareCamelCaseLambda implements Mustache.Lambda {
    private CodegenConfig generator = null;
    private Boolean escapeParam = false;
    private Boolean lowercaseFirstLetter = true;

    public SignareCamelCaseLambda(boolean lowercaseFirstLetter) {
        this.lowercaseFirstLetter = lowercaseFirstLetter;
    }

    public SignareCamelCaseLambda() {}

    public SignareCamelCaseLambda generator(final CodegenConfig generator) {
        this.generator = generator;
        return this;
    }

    public SignareCamelCaseLambda escapeAsParamName(final Boolean escape) {
        this.escapeParam = escape;
        return this;
    }

    @Override
    public void execute(Template.Fragment fragment, Writer writer) throws IOException {
        String text = camelize(fragment.execute().replace(" ", "_"), lowercaseFirstLetter);
        if (generator != null) {
            text = generator.sanitizeName(text);
            if (generator.reservedWords().contains(text)) {
                // Escaping must be done *after* camelize, because generators may escape using characters removed by camelize function.
                text = generator.escapeReservedWord(text);
            }

            if (escapeParam) {
                // NOTE: many generators call escapeReservedWord in toParamName, but we can't assume that's always the case.
                //       Here, we'll have to accept that we may be duplicating some work.
                text = generator.toParamName(text);
            }
        }
        writer.write(text);
    }


    /**
     * Taken partially from https://github.com/OpenAPITools/openapi-generator/blob/master/modules/openapi-generator/src/main/java/org/openapitools/codegen/utils/StringUtils.java
     *
     */

    public static final String NAME_CACHE_SIZE_PROPERTY = "org.openapitools.codegen.utils.namecache.cachesize";
    public static final String NAME_CACHE_EXPIRY_PROPERTY = "org.openapitools.codegen.utils.namecache.expireafter.seconds";


    static int cacheSize = Integer.parseInt(GlobalSettings.getProperty(NAME_CACHE_SIZE_PROPERTY, "200"));
    static int cacheExpiry = Integer.parseInt(GlobalSettings.getProperty(NAME_CACHE_EXPIRY_PROPERTY, "5"));
    static private Cache<Pair<String, Boolean>, String> camelizedWordsCache = Caffeine.newBuilder()
            .maximumSize(cacheSize)
                .expireAfterAccess(cacheExpiry, TimeUnit.SECONDS)
                .ticker(Ticker.systemTicker())
            .build();

    private static Pattern camelizeSlashPattern = Pattern.compile("\\/(.?)");
    private static Pattern camelizeUppercasePattern = Pattern.compile("(\\.?)(\\w)([^\\.]*)$");
    private static Pattern camelizeUnderscorePattern = Pattern.compile("(_)(.)");
    private static Pattern camelizeHyphenPattern = Pattern.compile("(-)(.)");
    private static Pattern camelizeDollarPattern = Pattern.compile("\\$");
    private static Pattern camelizeSimpleUnderscorePattern = Pattern.compile("_");

    /**
     * Camelize name (parameter, property, method, etc)
     *
     * @param inputWord                 string to be camelize
     * @param lowercaseFirstLetter lower case for first letter if set to true
     * @return camelized string
     */
    public static String camelize(final String inputWord, boolean lowercaseFirstLetter) {
        Pair<String, Boolean> key = new ImmutablePair<>(inputWord, lowercaseFirstLetter);

        return camelizedWordsCache.get(key, pair -> {
            String word = pair.getKey();
            Boolean lowerFirstLetter = pair.getValue();
            // Replace all slashes with dots (package separator)
            Matcher m = camelizeSlashPattern.matcher(word);
            while (m.find()) {
                word = m.replaceFirst("." + m.group(1)/*.toUpperCase()*/);
                m = camelizeSlashPattern.matcher(word);
            }

            // case out dots
            String[] parts = word.split("\\.");
            StringBuilder f = new StringBuilder();
            for (String z : parts) {
                if (z.length() > 0) {
                    f.append(Character.toUpperCase(z.charAt(0))).append(z.substring(1));
                }
            }
            word = f.toString();

            m = camelizeSlashPattern.matcher(word);
            while (m.find()) {
                word = m.replaceFirst(Character.toUpperCase(m.group(1).charAt(0)) + m.group(1).substring(1)/*.toUpperCase()*/);
                m = camelizeSlashPattern.matcher(word);
            }

            // Uppercase the class name.
            m = camelizeUppercasePattern.matcher(word);
            if (m.find()) {
                String rep = m.group(1) + m.group(2).toUpperCase(Locale.ROOT) + m.group(3);
                rep = camelizeDollarPattern.matcher(rep).replaceAll("\\\\\\$");
                word = m.replaceAll(rep);
            }

            // First turn to lowercase to fix SCREAMING_SNAKE_CASE
            word = word.toLowerCase(Locale.ROOT);

            // Remove all underscores (underscore_case to camelCase)
            m = camelizeUnderscorePattern.matcher(word);
            while (m.find()) {
                String original = m.group(2);
                String upperCase = original.toUpperCase(Locale.ROOT);
                if (original.equals(upperCase)) {
                    word = camelizeSimpleUnderscorePattern.matcher(word).replaceFirst("");
                } else {
                    word = m.replaceFirst(upperCase);
                }
                m = camelizeUnderscorePattern.matcher(word);
            }

            // Remove all hyphens (hyphen-case to camelCase)
            m = camelizeHyphenPattern.matcher(word);
            while (m.find()) {
                word = m.replaceFirst(m.group(2).toUpperCase(Locale.ROOT));
                m = camelizeHyphenPattern.matcher(word);
            }

            if (lowerFirstLetter && word.length() > 0) {
                int i = 0;
                char charAt = word.charAt(i);
                while (i + 1 < word.length() && !((charAt >= 'a' && charAt <= 'z') || (charAt >= 'A' && charAt <= 'Z'))) {
                    i = i + 1;
                    charAt = word.charAt(i);
                }
                i = i + 1;
                word = word.substring(0, i).toLowerCase(Locale.ROOT) + word.substring(i);
            }

            // remove all underscore
            word = camelizeSimpleUnderscorePattern.matcher(word).replaceAll("");
            return word;
        });
    }

}
